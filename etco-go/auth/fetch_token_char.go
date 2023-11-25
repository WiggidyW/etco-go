package auth

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"

	"github.com/lestrrat-go/jwx/jwk"
)

func tokenCharGet(
	x cache.Context,
	app esi.EsiApp,
	refreshToken string,
) (
	charId int32,
	expires time.Time,
	err error,
) {
	cacheKey := keys.CacheKeyTokenCharacter(uint8(app), refreshToken)
	return fetch.FetchWithCache(
		x,
		tokenCharGetFetchFunc(
			app,
			refreshToken,
			cacheKey,
		),
		cacheprefetch.TransientCache[int32](
			cacheKey,
			keys.TypeStrTokenCharacter,
			nil,
			nil,
		),
	)
}

func tokenCharGetFetchFunc(
	app esi.EsiApp,
	refreshToken string,
	cacheKey keys.Key,
) fetch.CachingFetch[int32] {
	return func(x cache.Context) (
		charId int32,
		expires time.Time,
		postFetch *cachepostfetch.Params,
		err error,
	) {
		x, cancel := x.WithCancel()
		defer cancel()

		// fetch JWKS in a goroutine
		chnJWKS := expirable.NewChanResult[jwk.Set](x.Ctx(), 1, 0)
		go expirable.P1Transceive(
			chnJWKS,
			x,
			esi.GetJWKS,
		)

		// fetch access token
		var accessToken string
		accessToken, expires, err = esi.GetAccessToken(x, refreshToken, app)
		if err != nil {
			return charId, expires, nil, err
		}

		// recv JWKS
		var jwks jwk.Set
		jwks, expires, err = chnJWKS.RecvExpMin(expires)
		if err != nil {
			return charId, expires, nil, err
		}

		// parse JWT
		charId, err = parseJWT(accessToken, jwks)
		if err != nil {
			return charId, expires, nil, err
		}

		postFetch = &cachepostfetch.Params{
			Set: cachepostfetch.LocalSetOne[int32](
				cacheKey,
				keys.TypeStrTokenCharacter,
				charId,
				expires,
			),
		}
		expires = fetch.CalcExpiresIn(expires, TOKEN_CHARACTER_MIN_EXPIRES)
		return charId, expires, postFetch, nil
	}
}
