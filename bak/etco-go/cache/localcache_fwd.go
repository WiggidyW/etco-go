package cache

import (
	"github.com/WiggidyW/etco-go/cache/localcache"
)

func RegisterType[T any](desc string, bufPoolCap int) string {
	return localcache.RegisterType[T](desc, bufPoolCap)
}

type BufferPool = localcache.BufferPool

func BufPool(typeStr string) *BufferPool {
	return localcache.BufPool(typeStr)
}
