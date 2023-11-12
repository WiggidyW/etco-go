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
	nsModified bool,
	rep *expirable.Expirable[REP],
	err error,
) {
	if params.CacheParams != nil {
		nsModified, rep, err = handleCache(x, *params.CacheParams)
		if nsModified || rep != nil || err != nil {
			return nsModified, rep, err
		}
	}
	return nsModified, rep, err
}
