package postfetch

import (
	"github.com/WiggidyW/etco-go/cache"
)

type Params struct {
	CacheParams *CacheParams
}

func Handle[REP any](
	x cache.Context,
	params *Params,
	fetchErr error,
) {
	var cacheParams *CacheParams = nil
	if params != nil {
		cacheParams = params.CacheParams
	}
	go handleCache(x, cacheParams, fetchErr)
}
