package bucket

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
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
	return fetch.FetchWithCache[map[string]struct{}](
		x,
		bundleKeysGetFetchFunc(getBuilder, cacheKey, typeStr),
		cacheprefetch.StrongCache[map[string]struct{}](
			cacheKey,
			typeStr,
			nil,
			nil,
		),
	)
}

func bundleKeysGetFetchFunc[V any](
	getBuilder fetch.HandledFetch[map[int32]map[string]V],
	cacheKey, typeStr string,
) fetch.CachingFetch[map[string]struct{}] {
	return func(x cache.Context) (
		rep map[string]struct{},
		expires time.Time,
		postFetch *cachepostfetch.Params,
		err error,
	) {
		var builder map[int32]map[string]V
		builder, expires, err = getBuilder(x)
		if err != nil {
			return nil, expires, nil, err
		}
		rep = extractBuilderBundleKeys(builder)
		postFetch = &cachepostfetch.Params{
			Set: cachepostfetch.ServerSetOne(cacheKey, typeStr, rep, expires),
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

func mapToKeySlice[K comparable, V any](
	m map[K]V,
) []K {
	keys := make([]K, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
