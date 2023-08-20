package cache

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/redis/go-redis/v9"
)

type SharedServerCache *redis.Client

type SharedClientCache *fastcache.Cache
