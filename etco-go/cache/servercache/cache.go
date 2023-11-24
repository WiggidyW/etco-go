package servercache

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
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
	key keys.Key,
) (
	val []byte,
	err error,
) {
	val, err = c.client.Get(ctx, key.String()).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			return nil, ErrServerGet{err}
		}
	}
	return val, nil
}

func (c Cache) set(
	ctx context.Context,
	key keys.Key,
	val []byte,
	ttl time.Duration,
) (err error) {
	err = c.client.Set(ctx, key.String(), val, ttl).Err()
	if err != nil {
		return ErrServerSet{err}
	}
	return nil
}

func (c Cache) del(
	ctx context.Context,
	key keys.Key,
) (err error) {
	err = c.client.Del(ctx, key.String()).Err()
	if err != nil {
		return ErrServerDel{err}
	}
	return nil
}
