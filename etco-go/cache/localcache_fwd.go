package cache

import (
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/cache/localcache"
)

func RegisterType[T any](desc string, bufPoolCap int) keys.Key {
	return localcache.RegisterType[T](desc, bufPoolCap)
}

type BufferPool = localcache.BufferPool

func BufPool(typeStr keys.Key) *BufferPool {
	return localcache.BufPool(typeStr)
}
