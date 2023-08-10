package authing

import (
	"context"
	"fmt"

	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/modelclient"
	"github.com/WiggidyW/weve-esi/util"
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
	useExtraIDs bool                    // whether to check alliance and corp IDs
	alrParams   AuthHashSetReaderParams // object name (domain key + access type)
	alrClient   AuthHashSetReaderClient // TODO: add override for skipping server cache to this type
	jwtClient   JWTCharacterClient
	charClient  *modelclient.ClientCharacterInfo
}

// Tries to return a refresh token in all cases if possible
func (ac *AuthingClient[F, D, C]) Fetch(
	ctx context.Context,
	params F,
) (*AuthingRep[D], error) {
	type rep = AuthingRep[D]
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the auth hash set in a separate goroutine
	chnHSSend, chnHSRecv := util.NewChanResult[AuthHashSet](ctx).Split()
	go ac.fetchHashSet(ctx, chnHSSend)

	// fetch a new token and the character ID from the provided token
	jwtRep, err := ac.jwtClient.Fetch(
		ctx,
		JWTCharacterParams{params.AuthRefreshToken()},
	)
	if err != nil {
		if jwtRep == nil {
			return nil, err
		} else { // return error with token
			return &rep{nil, false, jwtRep.RefreshToken}, err
		}
	}

	// if useExtraIDs is true, fetch character info in a separate goroutine
	var chnID util.ChanResult[modelclient.ModelCharacterInfo]
	if ac.useExtraIDs {
		chnID = util.NewChanResult[modelclient.ModelCharacterInfo](ctx)
		go ac.fetchExtraIDs(ctx, *jwtRep.CharacterID, chnID.ToSend())
	}

	// wait for the auth hash set
	hashSet, err := chnHSRecv.Recv()
	if err != nil { // return error with token
		return &rep{nil, false, jwtRep.RefreshToken}, err
	}

	// // check authorization

	if hashSet.ContainsCharacter(*jwtRep.CharacterID) { // character ID
		return ac.fetchAuthorized(ctx, params, jwtRep.RefreshToken)

	} else if ac.useExtraIDs { // extra IDs

		// wait for the extra IDs
		charInfo, err := chnID.Recv()
		if err != nil { // return error with token
			return &rep{nil, false, jwtRep.RefreshToken}, err
		}

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

	} // not authorized
	return &rep{nil, false, jwtRep.RefreshToken}, fmt.Errorf(
		"not authorized",
	)
}

func (ac *AuthingClient[F, D, C]) fetchHashSet(
	ctx context.Context,
	chnRes util.ChanSendResult[AuthHashSet],
) {
	if hashSet, err := ac.alrClient.Fetch(ctx, ac.alrParams); err != nil {
		chnRes.SendErr(err)
	} else {
		chnRes.SendOk(hashSet.Data())
	}
}

func (ac *AuthingClient[F, D, C]) fetchExtraIDs(
	ctx context.Context,
	characterId int32,
	chnRes util.ChanSendResult[modelclient.ModelCharacterInfo],
) {
	if charInfo, err := ac.charClient.Fetch(
		ctx,
		modelclient.NewFetchParamsCharacterInfo(characterId),
	); err != nil {
		chnRes.SendErr(err)
	} else {
		chnRes.SendOk(charInfo.Data())
	}
}

func (ac *AuthingClient[F, D, C]) fetchAuthorized(
	ctx context.Context,
	params F,
	token string,
) (*AuthingRep[D], error) {
	type rep = AuthingRep[D]
	if clientRep, err := ac.Client.Fetch(ctx, params); err != nil {
		return &rep{nil, true, token}, err
	} else {
		return &rep{clientRep, true, token}, nil
	}
}
