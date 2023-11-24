package market

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
)

type getEntries[E any] func(
	x cache.Context,
) (
	repOrStream esi.RepOrStream[E],
	expires time.Time,
	pages int,
	err error,
)

func getRaw[E marketOrdersEntry, M marketOrdersMap[E]](
	x cache.Context,
	getEntries getEntries[E],
	newMarketOrdersMap func() M,
	nsCacheKey, nsTypeStr keys.Key,
	cacheKey, typeStr keys.Key,
	minExpiresIn *time.Duration,
) (
	rep filteredMarketOrders,
	expires time.Time,
	err error,
) {
	return fetch.FetchWithCache[filteredMarketOrders](
		x,
		getRawFetchFunc[E, M](
			getEntries,
			newMarketOrdersMap,
			cacheKey,
			typeStr,
			nsCacheKey,
			nsTypeStr,
			minExpiresIn,
		),
		cacheprefetch.WeakMultiCacheDynamicKeys(
			cacheKey,
			typeStr,
			nsCacheKey,
			nsTypeStr,
			nil,
			cache.SloshTrue[filteredMarketOrders],
			nil,
		),
	)
}

func getRawFetchFunc[
	E marketOrdersEntry,
	M marketOrdersMap[E],
](
	getEntries getEntries[E],
	newMarketOrdersMap func() M,
	cacheKey, typeStr, nsCacheKey, nsTypeStr keys.Key,
	minExpiresIn *time.Duration,
) fetch.CachingFetch[filteredMarketOrders] {
	return func(x cache.Context) (
		rep filteredMarketOrders,
		expires time.Time,
		postFetch *cachepostfetch.Params,
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

		expires = fetch.CalcExpiresInOptional(expires, minExpiresIn)
		cacheSets := make([]cachepostfetch.ActionSet, 0, len(ordersWithKeys))

		var filtered filteredOrdersWithCacheKey
		for i := 0; i < len(ordersWithKeys); i++ {
			filtered = <-chn
			cacheSets = append(
				cacheSets,
				cachepostfetch.DualSet[filteredMarketOrders](
					filtered.CacheKey,
					typeStr,
					filtered.Orders,
					expires,
				),
			)
			if filtered.CacheKey.Bytes16() == cacheKey.Bytes16() {
				rep = filtered.Orders
			}
		}

		postFetch = &cachepostfetch.Params{
			Namespace: cachepostfetch.Namespace(
				nsCacheKey,
				nsTypeStr,
				expires,
			),
			Set: cacheSets,
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
