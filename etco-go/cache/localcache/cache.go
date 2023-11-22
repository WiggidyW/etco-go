package localcache

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/WiggidyW/etco-go/cache/keys"
)

type Cache struct {
	cache *fastcache.Cache // shared (1)
}

func newCache(maxBytes int) Cache {
	return Cache{cache: fastcache.New(maxBytes)}
}

func (c Cache) get(key keys.Key, dst []byte) []byte {
	val := c.cache.Get(dst, key.Bytes())
	if len(val) == 0 {
		return nil
	}
	return val
}

func (c Cache) del(key keys.Key) {
	c.cache.Del(key.Bytes())
}

func (c Cache) set(key keys.Key, val []byte) {
	c.cache.Set(key.Bytes(), val)
}
