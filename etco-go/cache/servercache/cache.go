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
	key [16]byte,
) (
	val []byte,
	err error,
) {
	k := string(key[:])
	val, err = c.client.Get(ctx, k).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			return nil, ErrServerGet{fmt.Errorf("%s: %w", k, err)}
		}
	}
	return val, nil
}

func (c Cache) set(
	ctx context.Context,
	key [16]byte,
	val []byte,
	ttl time.Duration,
) (err error) {
	k := string(key[:])
	err = c.client.Set(ctx, string(k), val, ttl).Err()
	if err != nil {
		return ErrServerSet{fmt.Errorf("%s: %w", k, err)}
	}
	return nil
}

func (c Cache) del(
	ctx context.Context,
	key [16]byte,
) (err error) {
	k := string(key[:])
	err = c.client.Del(ctx, k).Err()
	if err != nil {
		return ErrServerDel{fmt.Errorf("%s: %w", k, err)}
	}
	return nil
}
