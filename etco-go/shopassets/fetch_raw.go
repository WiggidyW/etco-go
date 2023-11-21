package shopassets

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
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
	return fetch.FetchWithCache(
		x,
		rawShopAssetsGetFetchFunc(locationId),
		cacheprefetch.WeakMultiCacheDynamicKeys(
			cacheKey,
			keys.TypeStrRawShopAssets,
			keys.CacheKeyNSRawShopAssets,
			keys.TypeStrNSRawShopAssets,
			nil,
			cache.SloshTrue[map[int32]int64],
			nil,
		),
	)
}

func rawShopAssetsGetFetchFunc(
	repLocationId int64,
) fetch.CachingFetch[map[int32]int64] {
	return func(x cache.Context) (
		assets map[int32]int64,
		expires time.Time,
		postFetch *cachepostfetch.Params,
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
		cacheSets := make([]cachepostfetch.ActionSet, 0, len(allRawAssets))

		for locationId, rawAssets := range allRawAssets {
			cacheSets = append(
				cacheSets,
				cachepostfetch.DualSet[map[int32]int64](
					keys.CacheKeyRawShopAssets(locationId),
					keys.TypeStrRawShopAssets,
					rawAssets,
					expires,
				),
			)
		}
		postFetch = &cachepostfetch.Params{
			Namespace: cachepostfetch.Namespace(
				keys.CacheKeyNSRawShopAssets,
				keys.TypeStrNSRawShopAssets,
				expires,
			),
			Set: cacheSets,
		}
		return allRawAssets[repLocationId], expires, nil, nil
	}
}
