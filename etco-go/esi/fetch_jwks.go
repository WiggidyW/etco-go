package esi

import (
	"encoding/json"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"

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
	var b *[]byte
	b, expires, err = fetch.HandleFetch[[]byte](
		x,
		&prefetch.Params[[]byte]{
			CacheParams: &prefetch.CacheParams[[]byte]{
				Get: prefetch.DualCacheGet(
					keys.CacheKeyJWKS, keys.TypeStrJWKS,
					true,
					newRep,
					cache.SloshTrue[[]byte],
				),
			},
		},
		jwksGetFetchFunc(newRep),
		EsiRetry,
	)
	if err != nil {
		return nil, expires, err
	}

	// unmarshal into a jwk.Set
	rep = jwk.NewSet()
	err = json.Unmarshal(*b, &rep)
	return rep, expires, err
}

func jwksGetNewRep(
	bufPool *cache.BufferPool,
) func() *[]byte {
	return func() *[]byte {
		b := make([]byte, 0, bufPool.Cap())
		return &b
	}
}

func jwksGetFetchFunc(
	newRep func() *[]byte,
) fetch.Fetch[[]byte] {
	return func(x cache.Context) (
		repPtr *[]byte,
		expires time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		var rep []byte
		rep, expires, err = getJWKS(x.Ctx(), newRep())
		if err != nil {
			return nil, expires, nil, err
		}
		postFetch = &postfetch.Params{
			CacheParams: &postfetch.CacheParams{
				Set: postfetch.DualCacheSetOne(
					keys.CacheKeyJWKS, keys.TypeStrJWKS,
					&rep,
					expires,
				),
			},
		}
		return &rep, expires, postFetch, nil
	}
}
