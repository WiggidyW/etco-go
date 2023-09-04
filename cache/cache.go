package cache

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/redis/go-redis/v9"
)

type SharedServerCache *redis.Client

func NewSharedServerCache(addr string) SharedServerCache {
	return redis.NewClient(&redis.Options{Addr: addr})
}

type SharedClientCache *fastcache.Cache

func NewSharedClientCache(maxBytes int) SharedClientCache {
	return fastcache.New(maxBytes)
}
