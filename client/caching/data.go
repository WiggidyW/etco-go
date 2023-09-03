package caching

import (
	"github.com/WiggidyW/etco-go/cache"
)

type CachingResponse[D any] struct {
	cache.ExpirableData[D]
	FromCache bool
}

func NewCachingResponse[D any](
	expirableData cache.ExpirableData[D],
	fromCache bool,
) *CachingResponse[D] {
	return &CachingResponse[D]{
		ExpirableData: expirableData,
		FromCache:     fromCache,
	}
}
