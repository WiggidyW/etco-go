package servercache

import (
	"context"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client     *redis.Client     // shared (1)
	lockClient *redislock.Client // shared (1)
}

func newCache(
	client *redis.Client,
) Cache {
	return Cache{
		client:     client,
		lockClient: redislock.New(*client),
	}
}

func (c Cache) get(
	ctx context.Context,
	key string,
) (
	val []byte,
	err error,
) {
	val, err = c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			return nil, ErrServerGet{fmt.Errorf("%s: %w", key, err)}
		}
	}
	return val, nil
}

func (c Cache) set(
	ctx context.Context,
	key string,
	val []byte,
	ttl time.Duration,
) (err error) {
	err = c.client.Set(ctx, key, val, ttl).Err()
	if err != nil {
		return ErrServerSet{fmt.Errorf("%s: %w", key, err)}
	}
	return nil
}

func (c Cache) del(
	ctx context.Context,
	key string,
) (err error) {
	err = c.client.Del(ctx, key).Err()
	if err != nil {
		return ErrServerDel{fmt.Errorf("%s: %w", key, err)}
	}
	return nil
}
