package auth

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/esi"
)

type ProtoLoginRep struct {
	RefreshToken string
	CharacterId  int32
	Admin        bool
}

func ProtoLogin(
	x cache.Context,
	app esi.EsiApp,
	accessCode string,
) (
	rep ProtoLoginRep,
	expires time.Time,
	err error,
) {
	rep.RefreshToken, err = esi.GetRefreshToken(x, accessCode, app)
	if err == nil {
		var infoRep ProtoTokenInfoRep
		infoRep, expires, err = ProtoTokenInfo(x, app, rep.RefreshToken)
		rep.CharacterId = infoRep.CharacterId
		rep.Admin = infoRep.Admin
	}
	return rep, expires, err
}

type ProtoTokenInfoRep struct {
	CharacterId int32
	Admin       bool
}

func ProtoTokenInfo(
	x cache.Context,
	app esi.EsiApp,
	refreshToken string,
) (
	rep ProtoTokenInfoRep,
	expires time.Time,
	err error,
) {
	if app != esi.EsiAppAuth {
		rep.CharacterId, expires, err =
			ProtoGetTokenCharacter(x, app, refreshToken)
	} else {
		var authRep AuthResponse
		authRep, expires, err = ProtoAdminAuthorized(x, refreshToken)
		rep.CharacterId = authRep.CharacterId
		rep.Admin = authRep.Authorized
	}
	return rep, expires, err
}
