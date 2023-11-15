package bucket

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

func bundleKeysGet[V any](
	x cache.Context,
	getBuilder fetch.HandledFetch[map[int32]map[string]V],
	cacheKey, typeStr string,
) (
	rep map[string]struct{},
	expires time.Time,
	err error,
) {
	return fetch.HandleFetch[map[string]struct{}](
		x,
		&prefetch.Params[map[string]struct{}]{
			CacheParams: &prefetch.CacheParams[map[string]struct{}]{
				Get: prefetch.ServerCacheGet[map[string]struct{}](
					cacheKey, typeStr,
					true,
					nil,
				),
			},
		},
		bundleKeysGetFetchFunc(getBuilder, cacheKey, typeStr),
		nil,
	)
}

func bundleKeysGetFetchFunc[V any](
	getBuilder fetch.HandledFetch[map[int32]map[string]V],
	cacheKey, typeStr string,
) fetch.Fetch[map[string]struct{}] {
	return func(x cache.Context) (
		rep map[string]struct{},
		expires time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		var builder map[int32]map[string]V
		builder, expires, err = getBuilder(x)
		if err != nil {
			return nil, expires, nil, err
		}
		rep = extractBuilderBundleKeys(builder)
		postFetch = &postfetch.Params{
			CacheParams: &postfetch.CacheParams{
				Set: postfetch.ServerCacheSetOne(
					cacheKey, typeStr,
					rep,
					expires,
				),
			},
		}
		return rep, expires, postFetch, nil
	}
}

func extractBuilderBundleKeys[V any](
	builder map[int32]map[string]V,
) map[string]struct{} {
	bundleKeys := make(map[string]struct{})
	for _, bundle := range builder {
		for bundleKey := range bundle {
			bundleKeys[bundleKey] = struct{}{}
		}
	}
	return bundleKeys
}
