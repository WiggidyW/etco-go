package internal

import (
	"context"

	ahs "github.com/WiggidyW/weve-esi/client/configure/authhashset"
	ahsreader "github.com/WiggidyW/weve-esi/client/configure/authhashset/reader"
	"github.com/WiggidyW/weve-esi/client/esi/jwt"
	charinfo "github.com/WiggidyW/weve-esi/client/esi/model/characterinfo"
	"github.com/WiggidyW/weve-esi/util"
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
	chnHSSend, chnHSRecv := util.NewChanResult[ahs.AuthHashSet](
		ctx,
	).Split()
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

	// if useExtraIDs is true, fetch IDs (in a separate goroutine)
	var chnID util.ChanResult[charinfo.CharacterInfoModel]
	if aci.useExtraIDs {
		chnID = util.NewChanResult[charinfo.CharacterInfoModel](ctx)
		go aci.fetchExtraIDs(ctx, *jwtRep.CharacterId, chnID.ToSend())
	}

	// wait for the auth hash set
	hashSet, err := chnHSRecv.Recv()
	if err != nil {
		return jwtRep, false, err
	}

	// check if characterID is authorized
	if hashSet.ContainsCharacter(*jwtRep.CharacterId) {
		return jwtRep, true, nil
	}

	// wait for the IDs and check if corpID or allianceID are authorized
	if aci.useExtraIDs {

		// wait for the IDs
		charInfo, err := chnID.Recv()
		if err != nil {
			return jwtRep, false, err
		}

		// check if corporationID or allianceID is authorized
		if (charInfo.AllianceId != nil &&
			hashSet.ContainsAlliance(*charInfo.AllianceId)) ||
			hashSet.ContainsCorporation(charInfo.CorporationId) {
			return jwtRep, true, nil
		}

	}

	// return not authorized
	return jwtRep, false, nil
}

func (aci AuthingClientInner) fetchHashSet(
	ctx context.Context,
	chnRes util.ChanSendResult[ahs.AuthHashSet],
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
	chnRes util.ChanSendResult[charinfo.CharacterInfoModel],
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
