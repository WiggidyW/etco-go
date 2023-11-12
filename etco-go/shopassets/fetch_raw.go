package shopassets

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

func rawShopAssetsGet(
	x cache.Context,
	locationId int64,
) (
	assets map[int32]int64,
	expires time.Time,
	err error,
) {
	cacheKey := keys.CacheKeyRawShopAssets(locationId)
	return fetch.HandleFetchVal(
		x,
		&prefetch.Params[map[int32]int64]{
			CacheParams: &prefetch.CacheParams[map[int32]int64]{
				Get: prefetch.DualCacheGet[map[int32]int64](
					cacheKey,
					keys.TypeStrRawShopAssets,
					true,
					nil,
					cache.SloshTrue[map[int32]int64],
				),
				Namespace: prefetch.CacheNamespace(
					keys.CacheKeyNSRawShopAssets,
					keys.TypeStrNSRawShopAssets,
					true,
				),
			},
		},
		rawShopAssetsGetFetchFunc(locationId),
		nil,
	)
}

func rawShopAssetsGetFetchFunc(
	repLocationId int64,
) fetch.Fetch[map[int32]int64] {
	return func(x cache.Context) (
		assets *map[int32]int64,
		expires time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		x, cancel := x.WithCancel()
		defer cancel()

		var repOrStream esi.RepOrStream[esi.AssetsEntry]
		var pages int
		repOrStream, expires, pages, err = esi.GetAssetsEntries(x)
		if err != nil {
			return nil, expires, nil, err
		}

		unflattenedAssets := newUnflattenedAssets(
			esi.ASSETS_ENTRIES_PER_PAGE * (pages - 1),
		)
		if repOrStream.Rep != nil {
			unflattenedAssets.addEntries(*repOrStream.Rep)
		} else /* if repOrStream.Stream != nil */ {
			var entries []esi.AssetsEntry
			for i := 0; i < pages; i++ {
				entries, expires, err = repOrStream.Stream.RecvExpMin(expires)
				if err != nil {
					return nil, expires, nil, err
				} else {
					unflattenedAssets.addEntries(entries)
				}
			}
		}

		allRawAssets := unflattenedAssets.flattenAndFilter()
		cacheSets := make([]postfetch.CacheActionSet, 0, len(allRawAssets))

		for locationId, rawAssets := range allRawAssets {
			cacheSets = append(cacheSets, postfetch.DualCacheSet(
				keys.CacheKeyRawShopAssets(locationId),
				keys.TypeStrRawShopAssets,
				&rawAssets,
				expires,
			))
		}
		postFetch = &postfetch.Params{
			CacheParams: &postfetch.CacheParams{
				Set: cacheSets,
				Namespace: postfetch.CacheNamespace(
					keys.CacheKeyNSRawShopAssets,
					keys.TypeStrNSRawShopAssets,
					expires,
				),
			},
		}
		return assets, expires, nil, nil
	}
}
