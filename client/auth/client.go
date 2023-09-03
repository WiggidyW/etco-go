package auth

import (
	"context"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/client/bucket"
	"github.com/WiggidyW/etco-go/client/esi/jwt"
	"github.com/WiggidyW/etco-go/client/esi/model/characterinfo"
)

type AuthClient struct {
	ahsReaderClient     bucket.SC_AuthHashSetReaderClient
	jwtClient           jwt.JWTClient
	characterInfoClient characterinfo.WC_CharacterInfoClient
}

func NewAuthClient(
	ahsReaderClient bucket.SC_AuthHashSetReaderClient,
	jwtClient jwt.JWTClient,
	characterInfoClient characterinfo.WC_CharacterInfoClient,
) AuthClient {
	return AuthClient{
		ahsReaderClient,
		jwtClient,
		characterInfoClient,
	}
}

// returns an error if not authorized, and returns the JWT response if possible
func (ac AuthClient) Fetch(
	ctx context.Context,
	params AuthParams,
) (
	authResponse AuthResponse,
	err error,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the auth hash set in a separate goroutine
	chnAuthHashSetSend, chnAuthHashSetRecv := chanresult.
		NewChanResult[b.AuthHashSet](ctx, 1, 0).Split()
	go ac.transceiveFetchAuthHashSet(
		ctx,
		params.AuthDomain,
		chnAuthHashSetSend,
	)

	// fetch a new token and the character ID from the provided token
	jwtRep, err := ac.fetchJWT(ctx, params.NativeRefreshToken)
	if err != nil {
		if jwtRep != nil {
			return AuthResponse{
				// Authorized: false,
				NativeRefreshToken: &jwtRep.NativeRefreshToken,
				CharacterId:        jwtRep.CharacterId,
				// CorporationId: nil,
				// AllianceId: nil,
			}, err
		} else {
			return AuthResponse{
				// Authorized: false,
				// NativeRefreshToken: nil,
				// CharacterId: nil,
				// CorporationId: nil,
				// AllianceId: nil,
			}, err
		}
	}

	// check if it's equal to the bootstrap admin character id
	// if so, no further checks are needed
	if *jwtRep.CharacterId == build.BOOTSTRAP_ADMIN_ID {
		return AuthResponse{
			Authorized:         true,
			NativeRefreshToken: &jwtRep.NativeRefreshToken,
			CharacterId:        jwtRep.CharacterId,
			// CorporationId: nil,
			// AllianceId: nil,
		}, nil
	}

	// fetch character info (in a separate goroutine)
	chnSendCharacterInfo, chnRecvCharacterInfo := chanresult.
		NewChanResult[characterinfo.CharacterInfoModel](ctx, 1, 0).
		Split()
	if params.UseExtraIds {
		go ac.transceiveFetchCharacterInfo(
			ctx,
			*jwtRep.CharacterId,
			chnSendCharacterInfo,
		)
	}

	// wait for the auth hash set
	authHashSet, err := chnAuthHashSetRecv.Recv()
	if err != nil {
		return AuthResponse{
			// Authorized: false,
			NativeRefreshToken: &jwtRep.NativeRefreshToken,
			CharacterId:        jwtRep.CharacterId,
			// CorporationId: nil,
			// AllianceId: nil,
		}, err
	}

	// check if characterID is authorized
	if authHashSet.PermittedCharacter(*jwtRep.CharacterId) {
		return AuthResponse{
			Authorized:         true,
			NativeRefreshToken: &jwtRep.NativeRefreshToken,
			CharacterId:        jwtRep.CharacterId,
			// CorporationId: nil,
			// AllianceId: nil,
		}, nil
	} else if authHashSet.BannedCharacter(*jwtRep.CharacterId) {
		return AuthResponse{
			// Authorized: false,
			NativeRefreshToken: &jwtRep.NativeRefreshToken,
			CharacterId:        jwtRep.CharacterId,
			// CorporationId: nil,
			// AllianceId: nil,
		}, nil
	}

	if !params.UseExtraIds ||
		(len(authHashSet.PermitAllianceIds) == 0 &&
			len(authHashSet.PermitCorporationIds) == 0) {
		// return not authorized
		return AuthResponse{
			// Authorized: false,
			NativeRefreshToken: &jwtRep.NativeRefreshToken,
			CharacterId:        jwtRep.CharacterId,
			// CorporationId: nil,
			// AllianceId: nil,
		}, nil
	}

	// wait for the IDs
	characterInfo, err := chnRecvCharacterInfo.Recv()
	if err != nil {
		return AuthResponse{
			// Authorized: false,
			NativeRefreshToken: &jwtRep.NativeRefreshToken,
			CharacterId:        jwtRep.CharacterId,
			// CorporationId: nil,
			// AllianceId: nil,
		}, err
	}

	// check if corporationID is authorized or banned
	if authHashSet.PermittedCorporation(
		characterInfo.CorporationId,
	) {
		return AuthResponse{
			Authorized:         true,
			NativeRefreshToken: &jwtRep.NativeRefreshToken,
			CharacterId:        jwtRep.CharacterId,
			CorporationId:      &characterInfo.CorporationId,
			AllianceId:         characterInfo.AllianceId,
		}, nil
	} else if authHashSet.BannedCorporation(
		characterInfo.CorporationId,
	) {
		return AuthResponse{
			Authorized:         false,
			NativeRefreshToken: &jwtRep.NativeRefreshToken,
			CharacterId:        jwtRep.CharacterId,
			CorporationId:      &characterInfo.CorporationId,
			AllianceId:         characterInfo.AllianceId,
		}, nil
	}

	// check if allianceID is authorized
	if characterInfo.AllianceId != nil &&
		authHashSet.PermittedAlliance(
			*characterInfo.AllianceId,
		) {
		return AuthResponse{
			Authorized:         true,
			NativeRefreshToken: &jwtRep.NativeRefreshToken,
			CharacterId:        jwtRep.CharacterId,
			CorporationId:      &characterInfo.CorporationId,
			AllianceId:         characterInfo.AllianceId,
		}, nil
	}

	// return not authorized
	return AuthResponse{
		// Authorized: false,
		NativeRefreshToken: &jwtRep.NativeRefreshToken,
		CharacterId:        jwtRep.CharacterId,
		CorporationId:      &characterInfo.CorporationId,
		AllianceId:         characterInfo.AllianceId,
	}, nil
}

// func (ac AuthClient) transceiveFetchJWT(
// 	ctx context.Context,
// 	nativeRefreshToken string,
// 	chnSend chanresult.ChanSendResult[jwt.JWTResponse],
// ) error {
// 	jwtRep, err := ac.fetchJWT(ctx, nativeRefreshToken)
// 	if err != nil {
// 		return chnSend.SendErr(err)
// 	} else {
// 		return chnSend.SendOk(jwtRep)
// 	}
// }

func (ac AuthClient) fetchJWT(
	ctx context.Context,
	nativeRefreshToken string,
) (
	jwtRep *jwt.JWTResponse,
	err error,
) {
	return ac.jwtClient.Fetch(
		ctx,
		jwt.JWTParams{NativeRefreshToken: nativeRefreshToken},
	)
}

func (ac AuthClient) transceiveFetchAuthHashSet(
	ctx context.Context,
	authDomain string,
	chnSend chanresult.ChanSendResult[b.AuthHashSet],
) error {
	authHashSet, err := ac.fetchAuthHashSet(ctx, authDomain)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(authHashSet)
	}
}

func (ac AuthClient) fetchAuthHashSet(
	ctx context.Context,
	authDomain string,
) (
	authHashSet b.AuthHashSet,
	err error,
) {
	authHashSetRep, err := ac.ahsReaderClient.Fetch(
		ctx,
		bucket.AuthHashSetReaderParams{AuthDomain: authDomain},
	)
	if err != nil {
		return authHashSet, err
	} else {
		return authHashSetRep.Data(), nil
	}
}

func (ac AuthClient) transceiveFetchCharacterInfo(
	ctx context.Context,
	characterId int32,
	chnSend chanresult.ChanSendResult[characterinfo.CharacterInfoModel],
) error {
	charInfo, err := ac.fetchCharacterInfo(ctx, characterId)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(charInfo)
	}
}

func (ac AuthClient) fetchCharacterInfo(
	ctx context.Context,
	characterId int32,
) (
	charInfo characterinfo.CharacterInfoModel,
	err error,
) {
	charInfoRep, err := ac.characterInfoClient.Fetch(
		ctx,
		characterinfo.CharacterInfoParams{CharacterId: characterId},
	)
	if err != nil {
		return charInfo, err
	} else {
		return charInfoRep.Data(), nil
	}
}
