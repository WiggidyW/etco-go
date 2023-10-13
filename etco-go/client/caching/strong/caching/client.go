package caching

import (
	"context"
	"fmt"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client"
	"github.com/WiggidyW/etco-go/client/caching"
	"github.com/WiggidyW/etco-go/logger"
)

type StrongCachingClient[
	F caching.CacheableParams, // the inner client params type
	D any, // the inner client response type
	ED cache.Expirable[D], // the inner client response type wrapped in an expirable
	C client.Client[F, ED], // the inner client type
] struct {
	Client     C
	cache      *cache.StrongCache[D, cache.ExpirableData[D]]
	minExpires time.Duration
}

func NewStrongCachingClient[
	F caching.CacheableParams,
	D any,
	ED cache.Expirable[D],
	C client.Client[F, ED],
](
	client C,
	minExpires time.Duration,
	sCache cache.SharedServerCache,
	serverLockTTL time.Duration,
	serverLockMaxWait time.Duration,
) StrongCachingClient[F, D, ED, C] {
	return StrongCachingClient[F, D, ED, C]{
		Client:     client,
		minExpires: minExpires,
		cache: cache.NewStrongCache[D, cache.ExpirableData[D]](
			sCache,
			serverLockTTL,
			serverLockMaxWait,
		),
	}
}

func (scc StrongCachingClient[F, D, ED, C]) GetAntiCache() *cache.StrongAntiCache {
	return scc.cache.ToAntiCache()
}

// func (scc StrongCachingClient[F, D, ED, C]) GetInnerClient() C {
// 	return scc.Client
// }

func (scc StrongCachingClient[F, D, ED, C]) Fetch(
	ctx context.Context,
	params F,
) (*caching.CachingResponse[D], error) {
	cacheKey := params.CacheKey()

	// lock the cache
	lock, err := scc.cache.Lock(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	// try to get from the cache
	if rep, err := scc.cache.Get(
		ctx,
		cacheKey,
		lock,
	); err != nil {
		// unlock in a goroutine and return the 'Get' error
		go func() { logger.Err(scc.cache.Unlock(lock)) }()
		return nil, err
	} else if rep != nil {
		// unlock blocking, and return an error if it fails
		if err := scc.cache.Unlock(lock); err != nil {
			return nil, err
		}
		// return the cached response
		logger.Info(fmt.Sprintf("%s: strong cache hit", cacheKey))
		return &caching.CachingResponse[D]{
			ExpirableData: *rep,
			FromCache:     true,
		}, nil
	}
	logger.Info(fmt.Sprintf("%s: strong cache miss", cacheKey))

	// fetch
	clientRep, err := scc.Client.Fetch(ctx, params)
	if err != nil {
		// unlock in a goroutine and return the 'Fetch' error
		go func() { logger.Err(scc.cache.Unlock(lock)) }()
		return nil, err
	}

	// initialize the new cache entry
	cacheEntry := caching.NewMinExpirableData[D, ED](
		*clientRep,
		scc.minExpires,
	)

	// cache the value
	if err := scc.cache.Set(
		cacheKey,
		cacheEntry,
		lock,
	); err != nil {
		// unlock in a goroutine and return the 'Set' error
		go func() { logger.Err(scc.cache.Unlock(lock)) }()
		return nil, err
	}

	// unlock blocking, and return an error if it fails
	if err := scc.cache.Unlock(lock); err != nil {
		return nil, err
	}

	return &caching.CachingResponse[D]{
		ExpirableData: cacheEntry,
		FromCache:     false,
	}, nil
}
