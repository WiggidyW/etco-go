package esi

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
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
	return fetch.FetchWithCache(
		x,
		func(x cache.Context) (
			accessToken string,
			expires time.Time,
			postFetch *cachepostfetch.Params,
			err error,
		) {
			accessToken, expires, err = fetch.FetchWithRetries(
				x,
				func(x cache.Context) (string, time.Time, error) {
					return authRefresh(x.Ctx(), refreshToken, app)
				},
				ESI_NUM_RETRIES,
				esiShouldRetry,
			)
			if err != nil {
				return accessToken, expires, nil, err
			}
			postFetch = &cachepostfetch.Params{
				Set: cachepostfetch.LocalSetOne[string](
					cacheKey,
					typeStr,
					accessToken,
					expires,
				),
			}
			return accessToken, expires, postFetch, nil
		},
		cacheprefetch.TransientCache[string](
			cacheKey,
			typeStr,
			nil,
			nil,
		),
	)
}
