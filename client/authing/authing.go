package authing

import (
	"context"
	"fmt"

	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/modelclient"
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

type AuthableParams interface {
	AuthRefreshToken() string // Native ESI refresh token
}

type AuthingClient[
	F AuthableParams, // the inner client's params type
	D any, // the inner client's response type
	C client.Client[F, D], // the inner client type
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
