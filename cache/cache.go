package cache

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/redis/go-redis/v9"
)

type SharedServerCache *redis.Client

func NewSharedServerCache() SharedServerCache {
	return &redis.Client{} // TODO
}

type SharedClientCache *fastcache.Cache

func NewSharedClientCache() SharedClientCache {
	return &fastcache.Cache{} // TODO
}
