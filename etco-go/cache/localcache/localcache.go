package localcache

import (
	"context"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
)

var (
	cache           Cache
	locksAndBufPool typeLocksAndBufPools
)

func init() {
	cache = newCache(build.CCACHE_MAX_BYTES)
	locksAndBufPool = newTypeLocksAndBufPools()
}

func BufPool(typeStr string) *BufferPool {
	return locksAndBufPool.get(typeStr).bufPool
}

func ObtainLock(
	ctx context.Context,
	key, typeStr string,
	maxWait time.Duration,
) (
	lock *Lock,
	err error,
) {
	return locksAndBufPool.get(typeStr).obtainLock(ctx, key, maxWait)
}

func RegisterType[T any](desc string, bufPoolCap int) string {
	typeStr, minBufPoolCap := NewTypeStr[T](desc)
	if minBufPoolCap > bufPoolCap {
		bufPoolCap = minBufPoolCap
	}
	locksAndBufPool.register(typeStr, bufPoolCap)
	return typeStr
}

func Get(key string, dst []byte) []byte {
	return cache.get(key, dst)
}

func Del(key string) {
	cache.del(key)
}

func Set(key string, val []byte) {
	cache.set(key, val)
}
