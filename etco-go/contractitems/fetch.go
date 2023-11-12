package contractitems

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

func contractItemsGet(
	x cache.Context,
	contractId int32,
) (
	rep []ContractItem,
	expires time.Time,
	err error,
) {
	cacheKey := keys.CacheKeyContractItems(contractId)
	return fetch.HandleFetchVal(
		x,
		&prefetch.Params[[]ContractItem]{
			CacheParams: &prefetch.CacheParams[[]ContractItem]{
				Get: prefetch.DualCacheGet[[]ContractItem](
					cacheKey, keys.TypeStrContractItems,
					true,
					contractItemsGetNewRep,
					cache.SloshTrue[[]ContractItem],
				),
			},
		},
		contractItemsGetFetchFunc(contractId, cacheKey),
		nil,
	)
}

func contractItemsGetNewRep() *[]ContractItem {
	rep := make([]ContractItem, 0, esi.CONTRACT_ITEMS_ENTRIES_MAKE_CAP)
	return &rep
}

func contractItemsGetFetchFunc(
	contractId int32,
	cacheKey string,
) fetch.Fetch[[]ContractItem] {
	return func(x cache.Context) (
		rep *[]ContractItem,
		expires time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		var entries []esi.ContractItemsEntry
		entries, expires, err = esi.GetContractItemsEntries(x, contractId)
		if err != nil {
			return nil, expires, nil, err
		}
		if entries != nil {
			rep = fromEntries(entries)
		}
		postFetch = &postfetch.Params{
			CacheParams: &postfetch.CacheParams{
				Set: postfetch.DualCacheSetOne(
					cacheKey,
					keys.TypeStrContractItems,
					rep,
					expires,
				),
			},
		}
		return rep, expires, postFetch, nil
	}
}
