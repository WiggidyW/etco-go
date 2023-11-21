package auth

import (
	"time"

	"github.com/WiggidyW/etco-go/bucket"
	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"

	// "github.com/WiggidyW/etco-go/proto"
	// "github.com/WiggidyW/etco-go/proto/protoerr"

	b "github.com/WiggidyW/etco-go-bucket"
)

type AdminStatus bool

const (
	IsAdmin AdminStatus = true
	Unknown AdminStatus = false
)

type AuthResponse struct {
	AdminStatus   AdminStatus // 'true' = isAdmin, 'false' = unknown
	Banned        bool
	Authorized    bool
	CharacterId   int32
	CorporationId *int32 // possibly nil if check wasn't needed
	AllianceId    *int32 // possibly nil if check wasn't needed or character not in alliance
}

func ProtoBoolUserAuthorized(
	x cache.Context,
	refreshToken string,
) (
	authorized bool,
	expires time.Time,
	err error,
) {
	var rep AuthResponse
	rep, expires, err = ProtoUserAuthorized(x, refreshToken)
	return rep.Authorized, expires, err
}

func ProtoBoolAdminAuthorized(
	x cache.Context,
	refreshToken string,
) (
	authorized bool,
	expires time.Time,
	err error,
) {
	var rep AuthResponse
	rep, expires, err = ProtoAdminAuthorized(x, refreshToken)
	return rep.Authorized, expires, err
}

func ProtoUserAuthorized(
	x cache.Context,
	refreshToken string,
) (
	rep AuthResponse,
	expires time.Time,
	err error,
) {
	rep, expires, err = protoAuthorized(x, refreshToken, "user", false)
	if !rep.Authorized &&
		err == nil &&
		!rep.Banned &&
		!(rep.CharacterId == 0) {
		var adminRep AuthResponse
		adminRep, expires, err = protoAuthorized(x, refreshToken, "admin", true)
		rep.AdminStatus = adminRep.AdminStatus
		rep.Authorized = adminRep.Authorized
		rep.Banned = adminRep.Banned
	}
	return rep, expires, err
}

func ProtoAdminAuthorized(
	x cache.Context,
	refreshToken string,
) (
	rep AuthResponse,
	expires time.Time,
	err error,
) {
	return protoAuthorized(x, refreshToken, "admin", true)
}

func protoAuthorized(
	x cache.Context,
	refreshToken string,
	domain string,
	isAdmin bool,
) (
	rep AuthResponse,
	expires time.Time,
	err error,
) {
	if refreshToken == "" {
		return rep, expires, nil
	}
	return fetch.HandleFetch(
		x,
		nil,
		func(x cache.Context) (
			AuthResponse,
			time.Time,
			*postfetch.Params,
			error,
		) {
			return protoFetchAuthorized(x, refreshToken, domain, isAdmin)
		},
		nil,
	)
}

func protoFetchAuthorized(
	x cache.Context,
	refreshToken string,
	domain string,
	isAdmin bool,
) (
	rep AuthResponse,
	expires time.Time,
	_ *postfetch.Params,
	err error,
) {
	x, cancel := x.WithCancel()
	defer cancel()

	if isAdmin {
		rep.AdminStatus = IsAdmin
		defer func() {
			if err != nil || !rep.Authorized {
				rep.AdminStatus = Unknown
			}
		}()
	}

	// fetch authhashset in a goroutine
	chnAHS := expirable.NewChanResult[b.AuthHashSet](x.Ctx(), 1, 0)
	go expirable.P2Transceive(
		chnAHS,
		x, domain,
		bucket.GetAuthHashSet,
	)

	// fetch characterID
	rep.CharacterId, expires, err = ProtoGetTokenCharacter(
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
	go expirable.P2Transceive(
		chnInfo,
		x, rep.CharacterId,
		esi.GetCharacterInfo,
	)

	// recv authhashset and check if character is permitted/banned
	var authHashSet b.AuthHashSet
	authHashSet, expires, err = chnAHS.RecvExpMin(expires)
	if err != nil {
		rep.Authorized = false
		return rep, expires, nil, err
	} else if authHashSet.PermittedCharacter(rep.CharacterId) {
		rep.Authorized = true
		return rep, expires, nil, nil
	} else if authHashSet.BannedCharacter(rep.CharacterId) {
		rep.Authorized = false
		rep.Banned = true
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
		rep.Banned = true
		return rep, expires, nil, nil
	} else if charInfo.AllianceId != nil && authHashSet.PermittedAlliance(*charInfo.AllianceId) {
		rep.Authorized = true
		return rep, expires, nil, nil
	} else {
		rep.Authorized = false
		return rep, expires, nil, nil
	}
}
