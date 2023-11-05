package bucket

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

func bundleKeysGet[V any](
	ctx context.Context,
	getBuilder fetch.HandledFetchVal[map[int32]map[string]V],
	typeStr, cacheKey string,
	lockTTL, lockMaxBackoff time.Duration,
) (
	rep map[string]struct{},
	expires *time.Time,
	err error,
) {
	var repPtr *map[string]struct{}
	repPtr, expires, err = fetch.HandleFetch[map[string]struct{}](
		ctx,
		&prefetch.Params[map[string]struct{}]{
			CacheParams: &prefetch.CacheParams[map[string]struct{}]{
				Get: prefetch.ServerCacheGet[map[string]struct{}](
					typeStr, cacheKey,
					lockTTL, lockMaxBackoff,
					nil,
				),
			},
		},
		bundleKeysGetFetchFunc(getBuilder, typeStr, cacheKey),
	)
	if err != nil {
		return nil, nil, err
	} else if repPtr != nil {
		rep = *repPtr
	}
	return rep, expires, nil
}

func bundleKeysGetFetchFunc[V any](
	getBuilder fetch.HandledFetchVal[map[int32]map[string]V],
	typeStr, cacheKey string,
) fetch.Fetch[map[string]struct{}] {
	return func(ctx context.Context) (
		rep *map[string]struct{},
		expires *time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		var builder map[int32]map[string]V
		builder, expires, err = getBuilder(ctx)
		if err != nil {
			return nil, nil, nil, err
		}
		rep = extractBuilderBundleKeys(builder)
		postFetch = &postfetch.Params{
			CacheParams: postfetch.ServerCacheSet(typeStr, cacheKey),
		}
		return rep, expires, postFetch, nil
	}
}

func extractBuilderBundleKeys[V any](
	builder map[int32]map[string]V,
) *map[string]struct{} {
	bundleKeys := make(map[string]struct{})
	for _, bundle := range builder {
		for bundleKey := range bundle {
			bundleKeys[bundleKey] = struct{}{}
		}
	}
	return &bundleKeys
}
