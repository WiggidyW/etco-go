package prefetch

import (
	"context"

	"github.com/WiggidyW/etco-go/cache/expirable"
)

type Params[REP any] struct {
	CacheParams *CacheParams[REP]
}

type UnhandledData struct {
	CacheLocks *CacheLocks
}

func Handle[REP any](
	ctx context.Context,
	params Params[REP],
) (
	rep *expirable.Expirable[REP],
	data *UnhandledData,
	err error,
) {
	if params.CacheParams != nil {
		data = &UnhandledData{}
		rep, data.CacheLocks, err = handleCache(ctx, *params.CacheParams)
		if err != nil {
			return nil, nil, err
		} else if rep != nil {
			data.CacheLocks = nil
		}
	}
	if data.CacheLocks == nil {
		data = nil
	}
	return rep, data, err
}
