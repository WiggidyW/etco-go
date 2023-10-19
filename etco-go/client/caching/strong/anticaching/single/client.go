package single

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	"github.com/WiggidyW/etco-go/client/caching"
	"github.com/WiggidyW/etco-go/logger"
)

// deletes cache entry after fetching
type StrongAntiCachingClient[
	F caching.AntiCacheableParams, // the inner client params type
	D any, // the inner client response type
	C client.Client[F, D], // the inner client type
] struct {
	Client    C
	antiCache *cache.StrongAntiCache
}

func NewStrongAntiCachingClient[
	F caching.AntiCacheableParams,
	D any,
	C client.Client[F, D],
](
	client C,
	antiCache *cache.StrongAntiCache,
) StrongAntiCachingClient[F, D, C] {
	return StrongAntiCachingClient[F, D, C]{
		Client:    client,
		antiCache: antiCache,
	}
}

func (sacc StrongAntiCachingClient[F, D, C]) InnerClient() C {
	return sacc.Client
}

func (sacc StrongAntiCachingClient[F, D, C]) Fetch(
	ctx context.Context,
	params F,
) (*D, error) {
	antiCacheKey := params.AntiCacheKey()

	if antiCacheKey == cachekeys.NULL_ANTI_CACHE_KEY {
		return sacc.Client.Fetch(ctx, params)
	}

	// lock the cache
	lock, err := sacc.antiCache.Lock(ctx, antiCacheKey)
	if err != nil {
		return nil, err
	}

	// cache delete
	if err := sacc.antiCache.Del(antiCacheKey, lock); err != nil {
		// unlock in a goroutine and return the 'Del' error
		go func() { logger.Err(sacc.antiCache.Unlock(lock)) }()
		return nil, err
	}

	// fetch
	rep, err := sacc.Client.Fetch(ctx, params)

	// the old cache entry has been deleted
	// the inner client has either written new ones or failed
	// thus, the locks have no further use.
	// (no reason to block the caller)
	go func() { logger.Err(sacc.antiCache.Unlock(lock)) }()

	// return the fetch result
	if err != nil {
		return nil, err
	} else {
		return rep, nil
	}
}

// func NewAntiCachingClient[
// 	F AntiCacheableParams,
// 	D any,
// 	CD any,
// 	C Client[F, D],
// ](
// 	client C,
// 	cache *cache.Cache[CD, cache.ExpirableData[CD]],
// ) AntiCachingClient[F, D, CD, C] {
// 	return AntiCachingClient[F, D, CD, C]{client, cache}
// }
