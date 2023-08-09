package authing

import (
	"context"
	"encoding/gob"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type CachingAuthHashSetReaderClient = *client.CachingClient[
	AuthHashSetReaderParams,
	AuthHashSet,
	cache.ExpirableData[AuthHashSet],
	authHashSetReaderClient,
]

type AuthHashSetReaderParams string // object name (domain key + access type)

func (rp AuthHashSetReaderParams) CacheKey() string {
	return string(rp)
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
	al, err := alrc.read(ctx, params.CacheKey())
	if err != nil {
		return nil, err
	}
	data := cache.NewExpirableData[AuthHashSet](
		*al,
		time.Now().Add(alrc.expires),
	)
	return &data, nil
}

type AntiCachingAuthHashSetWriterClient = *client.AntiCachingClient[
	AuthHashSetWriterParams,
	struct{},
	AuthHashSet,
	AuthHashSetWriterClient,
]

type AuthHashSetWriterParams struct {
	Key         string // object name (domain key + access type)
	AuthHashSet AuthHashSet
}

func (wp AuthHashSetWriterParams) AntiCacheKey() string {
	return wp.Key
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
