package servercache

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
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
	key keys.Key,
	ttl time.Duration, // this is a hard limit on how long we'll wait for the lock
	maxBackoff time.Duration, // if > ttl, it has no effect
	incrementBackoff time.Duration, // if > ttl, there'll only be 1 attempt
) (
	lock *Lock,
	err error,
) {
	var rawLock *redislock.Lock
	rawLock, err = cl.client.Obtain(
		ctx,
		key.String()+"lock",
		ttl,
		&redislock.Options{
			RetryStrategy: &incrementalRetry{
				currentBackoff:   incrementBackoff,
				maxBackoff:       maxBackoff,
				incrementBackoff: incrementBackoff,
			},
		},
	)
	if err != nil {
		err = ErrServerObtainLock{err}
	} else {
		lock = newLock(ctx, rawLock, ttl)
	}
	return lock, err
}

type incrementalRetry struct {
	currentBackoff   time.Duration
	maxBackoff       time.Duration
	incrementBackoff time.Duration
}

func (r *incrementalRetry) NextBackoff() time.Duration {
	if r.currentBackoff >= r.maxBackoff {
		return r.maxBackoff
	}
	backoff := r.currentBackoff
	r.currentBackoff += r.incrementBackoff
	return backoff
}
