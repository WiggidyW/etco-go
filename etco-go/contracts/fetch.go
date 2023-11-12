package contracts

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

func get(
	x cache.Context,
) (
	rep Contracts,
	expires time.Time,
	err error,
) {
	return fetch.HandleFetchVal(
		x,
		&prefetch.Params[Contracts]{
			CacheParams: &prefetch.CacheParams[Contracts]{
				Get: prefetch.DualCacheGet[Contracts](
					keys.CacheKeyContracts, keys.TypeStrContracts,
					true,
					nil,
					cache.SloshTrue[Contracts],
				),
			},
		},
		getFetchFunc,
		nil,
	)
}

func getFetchFunc(
	x cache.Context,
) (
	rep *Contracts,
	expires time.Time,
	postFetch *postfetch.Params,
	err error,
) {
	x, cancel := x.WithCancel()
	defer cancel()

	var repOrStream esi.RepOrStream[esi.ContractsEntry]
	var pages int
	repOrStream, expires, pages, err = esi.GetContractsEntries(x)
	if err != nil {
		return nil, expires, nil, err
	}

	rep = newContracts()
	if repOrStream.Rep != nil {
		rep.filterAddEntries(*repOrStream.Rep)
	} else /* if repOrStream.Stream != nil */ {
		var entries []esi.ContractsEntry
		for i := 0; i < pages; i++ {
			entries, expires, err = repOrStream.Stream.RecvExpMin(expires)
			if err != nil {
				return nil, expires, nil, err
			} else {
				rep.filterAddEntries(entries)
			}
		}
	}

	postFetch = &postfetch.Params{
		CacheParams: &postfetch.CacheParams{
			Set: postfetch.DualCacheSetOne(
				keys.CacheKeyContracts, keys.TypeStrContracts,
				rep,
				expires,
			),
		},
	}
	return rep, expires, postFetch, nil
}
