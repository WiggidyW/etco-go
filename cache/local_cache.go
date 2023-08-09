package cache

import (
	"sync"

	"github.com/VictoriaMetrics/fastcache"
)

type LocalCache struct {
	locks *sync.Map        // unique per type
	cache *fastcache.Cache // shared (1)
}

func newLocalCache(c *fastcache.Cache) *LocalCache {
	return &LocalCache{
		locks: new(sync.Map),
		cache: c,
	}
}

func (lc *LocalCache) lock(key string) *sync.Mutex {
	lockAny, _ := lc.locks.LoadOrStore(key, &sync.Mutex{})
	lock := lockAny.(*sync.Mutex)
	lock.Lock()
	return lock
}

func (lc *LocalCache) get(key string, dst []byte) []byte {
	val := lc.cache.Get(dst, []byte(key))
	if len(val) == 0 {
		return nil
	}
	return val
}

func (lc *LocalCache) del(key string) {
	lc.cache.Del([]byte(key))
}

func (lc *LocalCache) set(key string, val []byte) {
	lc.cache.Set([]byte(key), val)
}
