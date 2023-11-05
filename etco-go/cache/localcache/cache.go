package localcache

import (
	"github.com/VictoriaMetrics/fastcache"
)

type Cache struct {
	cache *fastcache.Cache // shared (1)
}

func newCache(maxBytes int) Cache {
	return Cache{cache: fastcache.New(maxBytes)}
}

func (c Cache) get(key string, dst []byte) []byte {
	val := c.cache.Get(dst, []byte(key))
	if len(val) == 0 {
		return nil
	}
	return val
}

func (c Cache) del(key string) {
	c.cache.Del([]byte(key))
}

func (c Cache) set(key string, val []byte) {
	c.cache.Set([]byte(key), val)
}
