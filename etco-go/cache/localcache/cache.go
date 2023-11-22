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

func (c Cache) get(key [16]byte, dst []byte) []byte {
	val := c.cache.Get(dst, key[:])
	if len(val) == 0 {
		return nil
	}
	return val
}

func (c Cache) del(key [16]byte) {
	c.cache.Del([]byte(key[:]))
}

func (c Cache) set(key [16]byte, val []byte) {
	c.cache.Set([]byte(key[:]), val)
}
