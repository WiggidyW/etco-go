package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

const INCREMENTAL_RETRY_INTERVAL = 30 * time.Millisecond

type ServerCache struct {
	client     *redis.Client     // shared (1)
	lockClient *redislock.Client // shared (1)
}

func newServerCache(
	client *redis.Client,
) *ServerCache {
	return &ServerCache{
		client:     client,
		lockClient: redislock.New(*client),
	}
}

func (sc *ServerCache) lock(
	ctx context.Context,
	key string,
	ttl time.Duration, // this is a hard limit on how long we'll wait for the lock
	maxBackoff time.Duration, // if > ttl, it has no effect
) (*redislock.Lock, error) {
	lock, err := sc.lockClient.Obtain(
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
		return nil, err
	}
	return lock, nil
}

func (sc *ServerCache) get(
	ctx context.Context,
	key string,
) ([]byte, error) {
	val, err := sc.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return val, nil
}

func (sc *ServerCache) set(
	ctx context.Context,
	key string,
	val []byte,
	ttl time.Duration,
) error {
	err := sc.client.Set(ctx, key, val, ttl).Err()
	if err != nil {
		return fmt.Errorf("error setting remote cache: %w", err)
	}
	return nil
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
