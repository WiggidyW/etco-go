package authing

import (
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/modelclient"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type AuthingRep[D any] struct {
	data         *D // nil if not authorized
	authorized   bool
	refreshToken string // new Native ESI refresh token
}

// panics if not authorized
func (ar *AuthingRep[D]) Data() D {
	return *ar.data
}

func (ar *AuthingRep[D]) Authorized() bool {
	return ar.authorized
}

func (ar *AuthingRep[D]) RefreshToken() string {
	return ar.refreshToken
}

type AuthableFetchParams interface {
	AuthRefreshToken() string // Native ESI refresh token
}

type AuthingClient[
	F AuthableFetchParams,
	D any,
	C client.Client[F, D],
] struct {
	Client      C
	useExtraIDs bool                           // whether to check alliance and corp IDs
	alrParams   AuthHashSetReaderParams        // object name (domain key + access type)
	alrClient   CachingAuthHashSetReaderClient // TODO: add override for skipping server cache to this type
	jwtClient   JWTCharacterClient
	charClient  *modelclient.ClientCharacterInfo
}

func (ac *AuthingClient[F, D, C]) notAuthorized(
	token string,
) (*AuthingRep[D], error) {
	return &AuthingRep[D]{
		refreshToken: token,
		authorized:   false,
		data:         nil,
	}, fmt.Errorf("not authorized")
}

func (ac *AuthingClient[F, D, C]) errWithToken(
	token string,
	err error,
) (*AuthingRep[D], error) {
	return &AuthingRep[D]{
		refreshToken: token,
		authorized:   false,
		data:         nil,
	}, err
}

func (ac *AuthingClient[F, D, C]) fetchAuthorized(
	ctx context.Context,
	params F,
	token string,
) (*AuthingRep[D], error) {
	rep, err := ac.Client.Fetch(ctx, params)
	if err != nil {
		return &AuthingRep[D]{
			refreshToken: token,
			authorized:   true,
			data:         nil,
		}, err
	} else {
		return &AuthingRep[D]{
			refreshToken: token,
			authorized:   true,
			data:         rep,
		}, nil
	}
}

// Tries to return a refresh token in all cases if possible
func (ac *AuthingClient[F, D, C]) Fetch(
	ctx context.Context,
	params F,
) (*AuthingRep[D], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	chnErr := make(chan error, 1)

	// fetch the auth hash set in a separate goroutine
	chnHashSet := make(chan *client.CachingRep[AuthHashSet], 1)
	go func() {
		if hashSet, err := ac.alrClient.Fetch(
			ctx,
			ac.alrParams,
		); err != nil {
			chnErr <- err
		} else {
			chnHashSet <- hashSet
		}
	}()

	// fetch a new token and the character ID from the provided token
	jwtRep, err := ac.jwtClient.Fetch(
		ctx,
		JWTCharacterParams{params.AuthRefreshToken()},
	)
	if err != nil {
		if jwtRep == nil {
			return nil, err
		} else {
			return ac.errWithToken(jwtRep.RefreshToken, err)
		}
	}

	// if useExtraIDs is true, fetch character info in a separate goroutine
	var chnCharInfo chan *client.CachingRep[modelclient.ModelCharacterInfo]
	if ac.useExtraIDs {
		chnCharInfo = make(
			chan *client.CachingRep[modelclient.ModelCharacterInfo],
			1,
		)
		go func() {
			if charInfo, err := ac.charClient.Fetch(
				ctx,
				modelclient.NewFetchParamsCharacterInfo(
					*jwtRep.CharacterID,
				),
			); err != nil {
				chnErr <- err
			} else {
				chnCharInfo <- charInfo
			}
		}()
	}

	// wait for the auth hash set
	var hashSet_ *client.CachingRep[AuthHashSet]
	select {
	case err := <-chnErr:
		return ac.errWithToken(jwtRep.RefreshToken, err)
	case hashSet_ = <-chnHashSet:
	}
	hashSet := hashSet_.Data()

	// check if the character ID is in the auth hash set
	if hashSet.ContainsCharacter(*jwtRep.CharacterID) {
		return ac.fetchAuthorized(ctx, params, jwtRep.RefreshToken)
	}

	// if useExtraIDs is true, check the corp and alliance IDs
	if ac.useExtraIDs {
		// wait for the character info
		var charInfo_ *client.CachingRep[modelclient.ModelCharacterInfo]
		select {
		case err := <-chnErr:
			return ac.errWithToken(jwtRep.RefreshToken, err)
		case charInfo_ = <-chnCharInfo:
		}
		charInfo := charInfo_.Data()

		// check if corporationID or allianceID is authorized
		if (charInfo.AllianceId != nil &&
			hashSet.ContainsAlliance(*charInfo.AllianceId)) ||
			hashSet.ContainsCorporation(charInfo.CorporationId) {
			return ac.fetchAuthorized(
				ctx,
				params,
				jwtRep.RefreshToken,
			)
		}
	}

	return ac.notAuthorized(jwtRep.RefreshToken)
}

type CachingAuthHashSetReaderClient = *client.CachingClient[
	AuthHashSetReaderParams,
	AuthHashSet,
	cache.ExpirableData[AuthHashSet],
	authHashSetReaderClient,
]

// object name (domain key + access type)
type AuthHashSetReaderParams string

func (alcf AuthHashSetReaderParams) Key() string {
	return string(alcf)
}

// Auth list reader client
type authHashSetReaderClient struct {
	*authHashSetClient
	expires time.Duration
}

// Gets the auth list and adds expiration
func (alrc authHashSetReaderClient) Fetch(
	ctx context.Context,
	params AuthHashSetReaderParams,
) (*cache.ExpirableData[AuthHashSet], error) {
	al, err := alrc.read(ctx, params.Key())
	if err != nil {
		return nil, err
	}
	data := cache.NewExpirableData[AuthHashSet](
		*al,
		time.Now().Add(alrc.expires),
	)
	return &data, nil
}

type AuthHashSetWriterParams struct {
	Key         string
	AuthHashSet AuthHashSet
}

// Auth list writer client
// singleton
type AuthHashSetWriterClient struct {
	*authHashSetClient
}

// Sets the auth list
func (aswc AuthHashSetWriterClient) Fetch(
	ctx context.Context,
	params AuthHashSetWriterParams,
) (*struct{}, error) {
	err := aswc.write(ctx, params.Key, params.AuthHashSet)
	if err != nil {
		return nil, err
	}
	return &struct{}{}, nil
}

type AuthHashSet struct {
	CharacterIDs   map[int32]struct{}
	CorporationIDs map[int32]struct{}
	AllianceIDs    map[int32]struct{}
}

func (ahs AuthHashSet) ContainsCharacter(id int32) bool {
	_, ok := ahs.CharacterIDs[id]
	return ok
}

func (ahs AuthHashSet) ContainsCorporation(id int32) bool {
	_, ok := ahs.CorporationIDs[id]
	return ok
}

func (ahs AuthHashSet) ContainsAlliance(id int32) bool {
	_, ok := ahs.AllianceIDs[id]
	return ok
}

// Basal auth list client
// singleton
// has methods to read and write auth lists
type authHashSetClient struct { // single one for the program
	_storageClient    *storage.Client
	storageClientOpts []option.ClientOption
	bucketName        string
	_bucketHandle     *storage.BucketHandle
	// bufPool *cache.BufferPool // same one used by local cache
}

// Gets the storage client (sets it if it's nil)
func (asc *authHashSetClient) storageClient() (*storage.Client, error) {
	if asc._storageClient == nil {
		ctx := context.Background()
		var err error
		if asc._storageClient, err = storage.NewClient(
			ctx,
			asc.storageClientOpts...,
		); err != nil {
			return nil, err
		}
	}
	return asc._storageClient, nil
}

// Gets the bucket handle (sets it if it's nil)
func (asc *authHashSetClient) bucketHandle() (*storage.BucketHandle, error) {
	if asc._bucketHandle == nil {
		storageClient, err := asc.storageClient()
		if err != nil {
			return nil, err
		}
		asc._bucketHandle = storageClient.Bucket(asc.bucketName)
	}
	return asc._bucketHandle, nil
}

// Gets the object handle from the provided key
func (asc *authHashSetClient) objHandle(
	key string,
) (*storage.ObjectHandle, error) {
	bucketHandle, err := asc.bucketHandle()
	if err != nil {
		return nil, err
	}
	return bucketHandle.Object(key), nil
}

// Reads an authHashSet from the provided key
func (asc *authHashSetClient) read(
	ctx context.Context,
	key string,
) (*AuthHashSet, error) {
	// get the object handle
	objHandle, err := asc.objHandle(key)
	if err != nil {
		return nil, err
	}
	// create the reader
	reader, err := objHandle.NewReader(ctx)
	if err != nil {
		// return empty if the object doesn't exist
		if err == storage.ErrObjectNotExist {
			return new(AuthHashSet), nil
		}
		return nil, err
	}
	defer reader.Close()
	// deserialize the auth list using Gob
	decoder := gob.NewDecoder(reader)
	authHashSet := new(AuthHashSet)
	err = decoder.Decode(authHashSet)
	if err != nil {
		return nil, err
	}

	return authHashSet, nil
}

// Writes the provided auth list to the provided key
func (asc *authHashSetClient) write(
	ctx context.Context,
	key string,
	obj AuthHashSet,
) error {
	// get the object handle
	objHandle, err := asc.objHandle(key)
	if err != nil {
		return err
	}
	// create the writer
	writer := objHandle.NewWriter(ctx)
	defer writer.Close()
	// serialize the auth list using Gob
	encoder := gob.NewEncoder(writer)
	err = encoder.Encode(obj)
	if err != nil {
		return err
	}

	return nil
}
