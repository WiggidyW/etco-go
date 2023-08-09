package client

import (
	"context"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
)

type CachingRep[D any] struct {
	cache.ExpirableData[D]
	fromCache bool
}

func (cr *CachingRep[D]) FromCache() bool {
	return cr.fromCache
}

func newCachingRep[D any](
	expirableData cache.ExpirableData[D],
	fromCache bool,
) *CachingRep[D] {
	return &CachingRep[D]{
		ExpirableData: expirableData,
		fromCache:     fromCache,
	}
}

type CacheableParams interface {
	CacheKey() string
}

// checks cache for data from the provided key before fetching
type CachingClient[
	F CacheableParams, // the inner client params type
	D any, // the inner client response type
	ED cache.Expirable[D], // the inner client response type wrapped in an expirable
	C Client[F, ED], // the inner client type
] struct {
	Client     C
	cache      *cache.Cache[D, cache.ExpirableData[D]]
	minExpires time.Duration
}

func NewCachingClient[
	F CacheableParams,
	D any,
	ED cache.Expirable[D],
	C Client[F, ED],
](
	client C,
	minExpires time.Duration,
	bufPool *cache.BufferPool,
	clientCache cache.SharedClientCache,
	serverCache cache.SharedServerCache,
	serverLockTTL time.Duration,
	serverLockMaxWait time.Duration,
) CachingClient[F, D, ED, C] {
	return CachingClient[F, D, ED, C]{
		Client:     client,
		minExpires: minExpires,
		cache: cache.NewCache[D, cache.ExpirableData[D]](
			bufPool,
			clientCache,
			serverCache,
			serverLockTTL,
			serverLockMaxWait,
		),
	}
}

func (cc *CachingClient[F, D, ED, C]) Fetch(
	ctx context.Context,
	params F,
) (*CachingRep[D], error) {
	cacheKey := params.CacheKey()

	// try to get from the cache
	rep, lock, err := cc.cache.GetOrLock(ctx, cacheKey)
	if err != nil {
		return nil, err
	} else if rep != nil { // cache hit
		return newCachingRep[D](*rep, true), nil
	} // cache miss

	// fetch
	edRepPtr, err := cc.Client.Fetch(ctx, params)
	if err != nil {
		cc.cache.Unlock(lock)
		return nil, err
	}
	edRep := *edRepPtr

	// set the expiry if needed
	edRepExpires := edRep.Expires()
	minExpiresTime := time.Now().Add(cc.minExpires)
	if edRepExpires.Before(minExpiresTime) {
		edRepExpires = minExpiresTime
	}

	// initialize the new cache entry
	cacheEntry := cache.NewExpirableData[D](edRep.Data(), edRepExpires)

	// cache the value
	if err := cc.cache.Set(
		ctx,
		cacheKey,
		cacheEntry,
		lock,
	); err != nil {
		return nil, err // doesn't return server errors
	}

	return newCachingRep(cacheEntry, false), nil
}

// func (cc *CachingClient[F, D, ED, C]) FetchDataOnly(
// 	ctx context.Context,
// 	params F,
// ) (D, error) {
// 	expirable, err := cc.FetchExpirableOnly(ctx, params)
// 	if err != nil {
// 		var d D
// 		return d, err
// 	}
// 	return expirable.Data, nil
// }
