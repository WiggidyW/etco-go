package servercache

import (
	"context"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/logger"

	"github.com/redis/go-redis/v9"
)

var (
	cache       Cache
	cacheLocker CacheLocker
)

func init() {
	rawClient := redis.NewClient(&redis.Options{Addr: build.SCACHE_ADDRESS})
	cache = newCache(rawClient)
	cacheLocker = newCacheLocker(rawClient)
}

func ObtainLock(
	ctx context.Context,
	key [16]byte,
	ttl, maxBackoff time.Duration,
) (*Lock, error) {
	return cacheLocker.lock(ctx, key, ttl, maxBackoff)
}

func Get(ctx context.Context, key [16]byte) ([]byte, error) {
	return cache.get(ctx, key)
}

func Set(ctx context.Context, key [16]byte, val []byte, ttl time.Duration) error {
	return cache.set(ctx, key, val, ttl)
}

func Del(ctx context.Context, key [16]byte) error {
	return cache.del(ctx, key)
}

func DelLogErr(key [16]byte) {
	err := Del(context.Background(), key)
	if err != nil {
		logger.Err(err.Error())
	}
}
