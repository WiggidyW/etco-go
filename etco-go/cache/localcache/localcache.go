package localcache

import (
	"context"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache/keys"
)

var (
	cache           Cache
	locksAndBufPool typeLocksAndBufPools
)

func init() {
	cache = newCache(build.CCACHE_MAX_BYTES)
	locksAndBufPool = newTypeLocksAndBufPools()
}

func BufPool(typeStr keys.Key) *BufferPool {
	return locksAndBufPool.get(typeStr).bufPool
}

func ObtainLock(
	ctx context.Context,
	key, typeStr keys.Key,
	maxWait time.Duration,
) (
	lock *Lock,
	err error,
) {
	return locksAndBufPool.get(typeStr).obtainLock(ctx, key, maxWait)
}

func RegisterType[T any](desc string, bufPoolCap int) keys.Key {
	typeStr, minBufPoolCap := NewTypeStr[T](desc)
	if minBufPoolCap > bufPoolCap {
		bufPoolCap = minBufPoolCap
	}
	locksAndBufPool.register(typeStr, bufPoolCap)
	return typeStr
}

func Get(key keys.Key, dst []byte) []byte {
	return cache.get(key, dst)
}

func Del(key keys.Key) {
	cache.del(key)
}

func Set(key keys.Key, val []byte) {
	cache.set(key, val)
}
