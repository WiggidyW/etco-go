package authorized

import (
	"time"

	"github.com/WiggidyW/etco-go/bucket"
	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"

	b "github.com/WiggidyW/etco-go-bucket"
)

type AuthResponse struct {
	Authorized    bool
	CharacterId   int32
	CorporationId *int32 // possibly nil if check wasn't needed
	AllianceId    *int32 // possibly nil if check wasn't needed or character not in alliance
}

func Authorized(
	x cache.Context,
	refreshToken string,
	domain string,
) (
	rep AuthResponse,
	expires time.Time,
	err error,
) {
	return fetch.HandleFetch(
		x,
		nil,
		func(x cache.Context) (
			AuthResponse,
			time.Time,
			*postfetch.Params,
			error,
		) {
			return fetchAuthorized(x, refreshToken, domain)
		},
		nil,
	)
}

func fetchAuthorized(
	x cache.Context,
	refreshToken string,
	domain string,
) (
	rep AuthResponse,
	expires time.Time,
	_ *postfetch.Params,
	err error,
) {
	x, cancel := x.WithCancel()
	defer cancel()

	// fetch authhashset in a goroutine
	chnAHS := expirable.NewChanResult[b.AuthHashSet](x.Ctx(), 1, 0)
	go expirable.Param2Transceive(
		chnAHS,
		x, domain,
		bucket.GetAuthHashSet,
	)

	// fetch characterID
	rep.CharacterId, expires, err = GetTokenCharacter(
		x,
		esi.EsiAppAuth,
		refreshToken,
	)
	if err != nil {
		return rep, expires, nil, err
	} else if rep.CharacterId == build.BOOTSTRAP_ADMIN_ID {
		rep.Authorized = true
		return rep, expires, nil, nil
	}

	// fetch character info in a goroutine
	chnInfo := expirable.NewChanResult[*esi.CharacterInfo](x.Ctx(), 1, 0)
	go expirable.Param2Transceive(
		chnInfo,
		x, rep.CharacterId,
		esi.GetCharacterInfo,
	)

	// recv authhashset and check if character is permitted/banned
	var authHashSet b.AuthHashSet
	authHashSet, expires, err = chnAHS.RecvExpMin(expires)
	if err != nil {
		return rep, expires, nil, err
	} else if authHashSet.PermittedCharacter(rep.CharacterId) {
		rep.Authorized = true
		return rep, expires, nil, nil
	} else if authHashSet.BannedCharacter(rep.CharacterId) {
		rep.Authorized = false
		return rep, expires, nil, nil
	}

	// recv characterinfo
	var charInfo *esi.CharacterInfo
	charInfo, expires, err = chnInfo.RecvExpMin(expires)
	if err != nil {
		return rep, expires, nil, err
	} else if charInfo == nil {
		rep.Authorized = false
		return rep, expires, nil, nil
	} else {
		rep.CorporationId = &charInfo.CorporationId
		rep.AllianceId = charInfo.AllianceId
	}

	// check if corp/alliance permitted/banned
	if authHashSet.PermittedCorporation(charInfo.CorporationId) {
		rep.Authorized = true
		return rep, expires, nil, nil
	} else if authHashSet.BannedCorporation(charInfo.CorporationId) {
		rep.Authorized = false
		return rep, expires, nil, nil
	} else if charInfo.AllianceId != nil && authHashSet.PermittedAlliance(*charInfo.AllianceId) {
		rep.Authorized = true
		return rep, expires, nil, nil
	} else {
		rep.Authorized = false
		return rep, expires, nil, nil
	}
}
