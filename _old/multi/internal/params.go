package internal

import "github.com/WiggidyW/weve-esi/cache"

type WeakDualCachingInnerParams[F any] struct {
	Params    F
	CacheKeys []string
	Locks     []*cache.WeakLock
}
