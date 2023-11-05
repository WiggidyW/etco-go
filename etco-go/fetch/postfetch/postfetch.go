package postfetch

import (
	"time"

	"github.com/WiggidyW/etco-go/fetch/prefetch"
	"github.com/WiggidyW/etco-go/logger"
)

type Params struct {
	CacheParams *CacheParams
}

func Handle[REP any](
	preFetchData *prefetch.UnhandledData,
	rep *REP,
	expires *time.Time,
	fetchErr error,
	params *Params,
) {
	if preFetchData != nil && preFetchData.CacheLocks != nil {
		var cacheParams *CacheParams = nil
		if params.CacheParams != nil {
			cacheParams = params.CacheParams
		}
		go func() {
			logger.MaybeErr(handleCache(
				*preFetchData.CacheLocks,
				cacheParams,
				expires,
				rep,
			))
		}()
	}
}
