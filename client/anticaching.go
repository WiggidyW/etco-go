package client

import (
	"context"

	"github.com/WiggidyW/weve-esi/cache"
)

type AntiCacheableParams interface {
	AntiCacheKey() string
}

// deletes cache entry after fetching
type AntiCachingClient[
	F AntiCacheableParams, // the inner client params type
	D any, // the inner client response type
	CD any, // the cache data type
	C Client[F, D], // the inner client type
] struct {
	Client C
	cache  *cache.Cache[CD, cache.ExpirableData[CD]]
}

func NewAntiCachingClient[
	F AntiCacheableParams,
	D any,
	CD any,
	C Client[F, D],
](
	client C,
	cache *cache.Cache[CD, cache.ExpirableData[CD]],
) AntiCachingClient[F, D, CD, C] {
	return AntiCachingClient[F, D, CD, C]{client, cache}
}

func (acc *AntiCachingClient[F, D, CD, C]) Fetch(
	ctx context.Context,
	params F,
) (*D, error) {
	antiCacheKey := params.AntiCacheKey()

	// cache lock
	lock := acc.cache.Lock(ctx, antiCacheKey)

	// fetch
	rep, err := acc.Client.Fetch(ctx, params)
	if err != nil {
		acc.cache.Unlock(lock)
		return nil, err
	}

	// cache delete
	acc.cache.Del(ctx, antiCacheKey, lock)

	return rep, nil
}
