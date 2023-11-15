package market

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

type getEntries[E any] func(
	x cache.Context,
) (
	repOrStream esi.RepOrStream[E],
	expires time.Time,
	pages int,
	err error,
)

func rawGet[E marketOrdersEntry, M marketOrdersMap[E]](
	x cache.Context,
	getEntries getEntries[E],
	newMarketOrdersMap func() M,
	nsCacheKey, nsTypeStr string,
	cacheKey, typeStr string,
	minExpiresIn *time.Duration,
) (
	rep filteredMarketOrders,
	expires time.Time,
	err error,
) {
	return fetch.HandleFetch[filteredMarketOrders](
		x,
		&prefetch.Params[filteredMarketOrders]{
			CacheParams: &prefetch.CacheParams[filteredMarketOrders]{
				Get: prefetch.DualCacheGet[filteredMarketOrders](
					cacheKey,
					typeStr,
					true,
					nil,
					cache.SloshTrue[filteredMarketOrders],
				),
				Namespace: prefetch.CacheNamespace(
					nsCacheKey,
					nsTypeStr,
					true,
				),
			},
		},
		rawGetFetchFunc[E, M](
			getEntries,
			newMarketOrdersMap,
			cacheKey,
			typeStr,
			nsCacheKey,
			nsTypeStr,
			minExpiresIn,
		),
		nil,
	)
}

func rawGetFetchFunc[
	E marketOrdersEntry,
	M marketOrdersMap[E],
](
	getEntries getEntries[E],
	newMarketOrdersMap func() M,
	cacheKey, typeStr, nsCacheKey, nsTypeStr string,
	minExpiresIn *time.Duration,
) fetch.Fetch[filteredMarketOrders] {
	return func(x cache.Context) (
		rep filteredMarketOrders,
		expires time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		x, cancel := x.WithCancel()
		defer cancel()

		var repOrStream esi.RepOrStream[E]
		var pages int
		repOrStream, expires, pages, err = getEntries(x)
		if err != nil {
			return rep, expires, nil, err
		}

		marketOrdersMap := newMarketOrdersMap()
		if repOrStream.Rep != nil {
			marketOrdersMapInsertEntries(marketOrdersMap, *repOrStream.Rep)
		} else {
			var entries []E
			for i := 0; i < pages; i++ {
				entries, expires, err = repOrStream.Stream.RecvExpMin(expires)
				if err != nil {
					return rep, expires, nil, err
				} else {
					marketOrdersMapInsertEntries(marketOrdersMap, entries)
				}
			}
		}

		ordersWithKeys := marketOrdersMap.GetAll(nsCacheKey)
		chn := make(chan filteredOrdersWithCacheKey, len(ordersWithKeys))
		for _, raw := range ordersWithKeys {
			go filterRawOrders(raw, chn)
		}

		expires = fetch.CalcExpiresOptional(expires, minExpiresIn)
		cacheSets := make([]postfetch.CacheActionSet, 0, len(ordersWithKeys))

		var filtered filteredOrdersWithCacheKey
		for i := 0; i < len(ordersWithKeys); i++ {
			filtered = <-chn
			cacheSets = append(cacheSets, postfetch.DualCacheSet(
				filtered.CacheKey,
				typeStr,
				filtered.Orders,
				expires,
			))
			if filtered.CacheKey == cacheKey {
				rep = filtered.Orders
			}
		}

		postFetch = &postfetch.Params{
			CacheParams: &postfetch.CacheParams{
				Set: cacheSets,
				Namespace: postfetch.CacheNamespace(
					nsCacheKey,
					nsTypeStr,
					expires,
				),
			},
		}
		return rep, expires, postFetch, nil
	}
}

func filterRawOrders(
	raw marketOrdersWithCacheKey,
	chn chan<- filteredOrdersWithCacheKey,
) {
	chn <- filteredOrdersWithCacheKey{
		CacheKey: raw.CacheKey,
		Orders:   newFilteredMarketOrders(*raw.Orders),
	}
}
