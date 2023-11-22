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

func BufPool(typeStr [16]byte) *BufferPool {
	return locksAndBufPool.get(typeStr).bufPool
}

func ObtainLock(
	ctx context.Context,
	key, typeStr [16]byte,
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
	locksAndBufPool.register(typeStr.Buf, bufPoolCap)
	return typeStr
}

func Get(key [16]byte, dst []byte) []byte {
	return cache.get(key, dst)
}

func Del(key [16]byte) {
	cache.del(key)
}

func Set(key [16]byte, val []byte) {
	cache.set(key, val)
}
