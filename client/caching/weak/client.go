package weak

// TODO: Move this to "single" directory (^)

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client"
	"github.com/WiggidyW/etco-go/client/caching"
	"github.com/WiggidyW/etco-go/logger"
)

// TODO: (Critical) only use server cache, only use client cache, or use both (OPTION PARAMETERS)
// checks cache for data from the provided key before fetching
type WeakCachingClient[
	F caching.CacheableParams, // the inner client params type
	D any, // the inner client response type
	ED cache.Expirable[D], // the inner client response type wrapped in an expirable
	C client.Client[F, ED], // the inner client type
] struct {
	Client     C
	cache      *cache.WeakCache[D, cache.ExpirableData[D]]
	minExpires time.Duration
}

func NewWeakCachingClient[
	F caching.CacheableParams,
	D any,
	ED cache.Expirable[D],
	C client.Client[F, ED],
](
	client C,
	minExpires time.Duration,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
	serverLockTTL time.Duration,
	serverLockMaxWait time.Duration,
) WeakCachingClient[F, D, ED, C] {
	return WeakCachingClient[F, D, ED, C]{
		Client:     client,
		minExpires: minExpires,
		cache: cache.NewWeakCache[D, cache.ExpirableData[D]](
			cCache,
			sCache,
			serverLockTTL,
			serverLockMaxWait,
		),
	}
}

func (wcc WeakCachingClient[F, D, ED, C]) Fetch(
	ctx context.Context,
	params F,
) (*caching.CachingResponse[D], error) {
	cacheKey := params.CacheKey()

	// try to get from the cache
	cacheRep, lock := wcc.cache.GetOrLock(ctx, cacheKey)

	// return now if it was a cache hit
	if cacheRep != nil {
		return &caching.CachingResponse[D]{
			ExpirableData: *cacheRep,
			FromCache:     true,
		}, nil
	}

	// fetch
	clientRep, err := wcc.Client.Fetch(ctx, params)
	if err != nil {
		go func() { logger.Err(wcc.cache.Unlock(lock)) }()
		return nil, err
	}

	// initialize the new cache entry
	cacheEntry := caching.NewMinExpirableData[D, ED](
		*clientRep,
		wcc.minExpires,
	)

	// cache the value in the background, logging any errors
	go func() {
		logger.Err(wcc.cache.Set(
			cacheKey,
			cacheEntry,
			lock,
		))
	}()

	return &caching.CachingResponse[D]{
		ExpirableData: cacheEntry,
		FromCache:     false,
	}, nil
}
