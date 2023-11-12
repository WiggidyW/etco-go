package esi

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

func accessTokenGet(
	x cache.Context,
	refreshToken string,
	app EsiApp,
) (
	accessToken string,
	expires time.Time,
	err error,
) {
	typeStr := app.TypeStrToken()
	cacheKey := app.CacheKeyToken(refreshToken)
	return fetch.HandleFetchVal(
		x,
		&prefetch.Params[string]{
			CacheParams: &prefetch.CacheParams[string]{
				Get: prefetch.LocalCacheGet[string](
					app.CacheKeyToken(refreshToken),
					app.TypeStrToken(),
					true,
					nil,
				),
			},
		},
		accessTokenGetFetchFunc(refreshToken, app, cacheKey, typeStr),
		EsiRetry,
	)
}

func accessTokenGetFetchFunc(
	refreshToken string,
	app EsiApp,
	cacheKey, typeStr string,
) fetch.Fetch[string] {
	return func(x cache.Context) (
		accessTokenPtr *string,
		expires time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		var accessToken string
		accessToken, expires, err = authRefresh(x.Ctx(), refreshToken, app)
		if err != nil {
			return nil, expires, nil, err
		}
		postFetch = &postfetch.Params{
			CacheParams: &postfetch.CacheParams{
				Set: postfetch.LocalCacheSetOne(
					cacheKey, typeStr,
					&accessToken,
					expires,
				),
			},
		}
		return &accessToken, expires, postFetch, nil
	}
}
