package internal

import (
	"context"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	build "github.com/WiggidyW/etco-go/buildconstants"
	ahsreader "github.com/WiggidyW/etco-go/client/configure/authhashset/reader"
	"github.com/WiggidyW/etco-go/client/esi/jwt"
	charinfo "github.com/WiggidyW/etco-go/client/esi/model/characterinfo"
)

type AuthingClientInner struct {
	useExtraIDs bool                              // whether to check alliance and corp IDs
	alrParams   ahsreader.AuthHashSetReaderParams // object name (domain key + access type)
	alrClient   ahsreader.AuthHashSetReaderClient // TODO: add override for skipping server cache to this type
	jwtClient   jwt.JWTClient
	charClient  charinfo.WC_CharacterInfoClient
}

// returns an error if not authorized, and returns the JWT response if possible
func (aci AuthingClientInner) Authorized(
	ctx context.Context,
	nativeRefreshToken string, // character
) (_ *jwt.JWTResponse, authorized bool, _ error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the auth hash set in a separate goroutine
	chnHSSend, chnHSRecv := chanresult.
		NewChanResult[b.AuthHashSet](ctx, 0, 0).Split()
	go aci.fetchHashSet(ctx, chnHSSend)

	// fetch a new token and the character ID from the provided token
	jwtRep, err := aci.jwtClient.Fetch(
		ctx,
		jwt.JWTParams{NativeRefreshToken: nativeRefreshToken},
	)
	if err != nil {
		if jwtRep == nil {
			return nil, false, err
		} else {
			return jwtRep, false, err
		}
	}

	// check if it's equal to the bootstrap admin character id
	// if so, no further checks are needed
	if *jwtRep.CharacterId == build.BOOTSTRAP_ADMIN_ID {
		return jwtRep, true, nil
	}

	// if useExtraIDs is true, fetch IDs (in a separate goroutine)
	var chnID chanresult.ChanResult[charinfo.CharacterInfoModel]
	if aci.useExtraIDs {
		chnID = chanresult.
			NewChanResult[charinfo.CharacterInfoModel](ctx, 0, 0)
		go aci.fetchExtraIDs(ctx, *jwtRep.CharacterId, chnID.ToSend())
	}

	// wait for the auth hash set
	authHashSet, err := chnHSRecv.Recv()
	if err != nil {
		return jwtRep, false, err
	}

	// check if characterID is authorized
	if authHashSet.PermittedCharacter(*jwtRep.CharacterId) {
		return jwtRep, true, nil
	} else if authHashSet.BannedCharacter(*jwtRep.CharacterId) {
		return jwtRep, false, nil
	}

	// wait for the IDs and check if corpID or allianceID are authorized
	if aci.useExtraIDs && (len(authHashSet.PermitAllianceIds) > 0 ||
		len(authHashSet.PermitCorporationIds) > 0) {

		// wait for the IDs
		charInfo, err := chnID.Recv()
		if err != nil {
			return jwtRep, false, err
		}

		// check if corporationID is authorized or banned
		if authHashSet.PermittedCorporation(charInfo.CorporationId) {
			return jwtRep, true, nil
		} else if authHashSet.BannedCorporation(charInfo.CorporationId) {
			return jwtRep, false, nil
		}

		// check if allianceID is authorized
		if charInfo.AllianceId != nil &&
			authHashSet.PermittedAlliance(*charInfo.AllianceId) {
			if authHashSet.PermittedAlliance(*charInfo.AllianceId) {
				return jwtRep, true, nil
			}
		}
	}

	// return not authorized
	return jwtRep, false, nil
}

func (aci AuthingClientInner) fetchHashSet(
	ctx context.Context,
	chnRes chanresult.ChanSendResult[b.AuthHashSet],
) {
	if hashSet, err := aci.alrClient.Fetch(
		ctx,
		aci.alrParams,
	); err != nil {
		chnRes.SendErr(err)
	} else {
		chnRes.SendOk(hashSet.Data())
	}
}

func (aci AuthingClientInner) fetchExtraIDs(
	ctx context.Context,
	characterId int32,
	chnRes chanresult.ChanSendResult[charinfo.CharacterInfoModel],
) {
	if charInfo, err := aci.charClient.Fetch(
		ctx,
		charinfo.CharacterInfoParams{CharacterId: characterId},
	); err != nil {
		chnRes.SendErr(err)
	} else {
		chnRes.SendOk(charInfo.Data())
	}
}
