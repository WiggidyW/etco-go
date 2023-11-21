package esi

import (
	"encoding/json"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"

	"github.com/lestrrat-go/jwx/jwk"
)

func jwksGet(
	x cache.Context,
) (
	rep jwk.Set,
	expires time.Time,
	err error,
) {
	bufPool := cache.BufPool(keys.TypeStrJWKS)
	newRep := jwksGetNewRep(bufPool)

	// fetch JWKS bytes
	var b []byte
	b, expires, err = fetch.FetchWithCache[[]byte](
		x,
		jwksGetFetchFunc(newRep),
		cacheprefetch.WeakCache(
			keys.CacheKeyJWKS,
			keys.TypeStrJWKS,
			newRep,
			cache.SloshTrue[[]byte],
			nil,
		),
	)
	if err != nil {
		return nil, expires, err
	}

	// unmarshal into a jwk.Set
	rep = jwk.NewSet()
	err = json.Unmarshal(b, &rep)
	return rep, expires, err
}

func jwksGetNewRep(
	bufPool *cache.BufferPool,
) func() []byte {
	return func() []byte {
		return make([]byte, 0, bufPool.Cap())
	}
}

func jwksGetFetchFunc(newRep func() []byte) fetch.CachingFetch[[]byte] {
	return func(x cache.Context) (
		rep []byte,
		expires time.Time,
		postFetch *cachepostfetch.Params,
		err error,
	) {
		rep, expires, err = fetch.FetchWithRetries(
			x,
			func(x cache.Context) (rep []byte, expires time.Time, err error) {
				return getJWKS(x.Ctx(), newRep())
			},
			ESI_NUM_RETRIES,
			esiShouldRetry,
		)
		if err != nil {
			return nil, expires, nil, err
		}
		postFetch = &cachepostfetch.Params{
			Set: cachepostfetch.DualSetOne[[]byte](
				keys.CacheKeyJWKS,
				keys.TypeStrJWKS,
				rep,
				expires,
			),
		}
		return rep, expires, postFetch, nil
	}
}
