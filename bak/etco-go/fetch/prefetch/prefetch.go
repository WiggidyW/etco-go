package prefetch

import (
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
)

type Params[REP any] struct {
	CacheParams *CacheParams[REP]
}

func Handle[REP any](
	x cache.Context,
	params Params[REP],
) (
	ncRetry bool,
	rep *expirable.Expirable[REP],
	err error,
) {
	if params.CacheParams != nil {
		ncRetry, rep, err = handleCache(x, *params.CacheParams)
		if ncRetry || rep != nil || err != nil {
			return ncRetry, rep, err
		}
	}
	return ncRetry, rep, err
}
