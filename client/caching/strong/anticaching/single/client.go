package single

import (
	"context"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/caching"
	"github.com/WiggidyW/weve-esi/logger"
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

func (sacc StrongAntiCachingClient[F, D, C]) Fetch(
	ctx context.Context,
	params F,
) (*D, error) {
	antiCacheKey := params.AntiCacheKey()

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
