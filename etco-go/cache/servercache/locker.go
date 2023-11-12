package servercache

import (
	"context"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type CacheLocker struct {
	client *redislock.Client
}

func newCacheLocker(
	client *redis.Client,
) CacheLocker {
	return CacheLocker{client: redislock.New(*client)}
}

func (cl CacheLocker) lock(
	ctx context.Context,
	key string,
	ttl time.Duration, // this is a hard limit on how long we'll wait for the lock
	maxBackoff time.Duration, // if > ttl, it has no effect
) (
	lock *Lock,
	err error,
) {
	var rawLock *redislock.Lock
	rawLock, err = cl.client.Obtain(
		ctx,
		key+"lock",
		ttl,
		&redislock.Options{
			RetryStrategy: &incrementalRetry{
				maxBackoff: maxBackoff,
			},
		},
	)
	if err != nil {
		err = ErrServerObtainLock{fmt.Errorf("%s: %w", key, err)}
	} else {
		lock = newLock(ctx, rawLock, ttl)
	}
	return lock, err
}

type incrementalRetry struct {
	currentBackoff time.Duration
	maxBackoff     time.Duration
}

func (r *incrementalRetry) NextBackoff() time.Duration {
	if r.currentBackoff >= r.maxBackoff {
		return 0
	}
	backoff := r.currentBackoff
	r.currentBackoff += INCREMENTAL_RETRY_INTERVAL
	return backoff
}
