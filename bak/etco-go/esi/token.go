package esi

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
)

const (
	TOKEN_BUF_CAP int = 0
)

func init() {
	keys.TypeStrAuthToken = cache.RegisterType[string](EsiAppAuth.String(), TOKEN_BUF_CAP)
	keys.TypeStrCorpToken = cache.RegisterType[string](EsiAppCorp.String(), TOKEN_BUF_CAP)
	keys.TypeStrStructureInfoToken = cache.RegisterType[string](EsiAppStructureInfo.String(), TOKEN_BUF_CAP)
	keys.TypeStrMarketsToken = cache.RegisterType[string](EsiAppMarkets.String(), TOKEN_BUF_CAP)
}

func GetAccessToken(
	x cache.Context,
	refreshToken string,
	app EsiApp,
) (
	accessToken string,
	expires time.Time,
	err error,
) {
	return accessTokenGet(x, refreshToken, app)
}
