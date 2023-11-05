package localcache

import (
	build "github.com/WiggidyW/etco-go/buildconstants"
)

var (
	cache           Cache
	locksAndBufPool TypeLocksAndBufPools
)

func init() {
	cache = newCache(build.CCACHE_MAX_BYTES)
	locksAndBufPool = newTypeLocksAndBufPools()
}

func GetLocksAndBufPool(typeStr string) TypeLocksAndBufPool {
	return locksAndBufPool.get(typeStr)
}

func RegisterType[T any](bufPoolCap int) string {
	typeStr := NewTypeStr[T]()
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
