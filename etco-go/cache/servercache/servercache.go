package servercache

import (
	"context"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache/keys"
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
	key keys.Key,
	ttl, maxBackoff, incrementBackoff time.Duration,
) (*Lock, error) {
	return cacheLocker.lock(ctx, key, ttl, maxBackoff, incrementBackoff)
}

func Get(ctx context.Context, key keys.Key) ([]byte, error) {
	return cache.get(ctx, key)
}

func Set(ctx context.Context, key keys.Key, val []byte, ttl time.Duration) error {
	return cache.set(ctx, key, val, ttl)
}

func Del(ctx context.Context, key keys.Key) error {
	return cache.del(ctx, key)
}

func DelLogErr(key keys.Key) {
	err := Del(context.Background(), key)
	if err != nil {
		logger.Err(err.Error())
	}
}
