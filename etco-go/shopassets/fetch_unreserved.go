package shopassets

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
	"github.com/WiggidyW/etco-go/purchasequeue"
)

func unreservedShopAssetsGet(
	x cache.Context,
	locationId int64,
) (
	assets map[int32]int64,
	expires time.Time,
	err error,
) {
	cacheKey := keys.CacheKeyUnreservedShopAssets(locationId)
	return fetch.FetchWithCache(
		x,
		unreservedShopAssetsGetFetchFunc(cacheKey, locationId),
		cacheprefetch.StrongCache[map[int32]int64](
			cacheKey,
			keys.TypeStrUnreservedShopAssets,
			nil,
			nil,
		),
	)
}

func unreservedShopAssetsGetFetchFunc(
	cacheKey keys.Key,
	locationId int64,
) fetch.CachingFetch[map[int32]int64] {
	return func(x cache.Context) (
		assets map[int32]int64,
		expires time.Time,
		postFetch *cachepostfetch.Params,
		err error,
	) {
		// Cancel goroutines if we return early
		x, cancel := x.WithCancel()
		defer cancel()

		// Get the location's full assets in a goroutine
		chnAssets := expirable.NewChanResult[map[int32]int64](x.Ctx(), 1, 0)
		go expirable.P2Transceive(
			chnAssets,
			x, locationId,
			getRawShopAssets,
		)

		// Get the location's purchase queue items
		var reservedItems map[int32]int64
		reservedItems, expires, err =
			purchasequeue.GetLocationPurchaseQueueItems(x, locationId)
		if err != nil {
			return nil, expires, nil, err
		}

		// Recv the location's full assets
		assets, expires, err = chnAssets.RecvExpMin(expires)
		if err != nil {
			return nil, expires, nil, err
		}

		// Subtract the reserved items from the full assets
		var oldQuantity, newQuantity int64
		var ok bool
		for typeId, quantity := range reservedItems {
			oldQuantity, ok = assets[typeId]
			if ok {
				newQuantity = oldQuantity - quantity
				if newQuantity > 0 {
					assets[typeId] = newQuantity
				} else {
					delete(assets, typeId)
				}
			}
		}

		postFetch = &cachepostfetch.Params{
			Set: cachepostfetch.ServerSetOne[map[int32]int64](
				cacheKey,
				keys.TypeStrUnreservedShopAssets,
				assets,
				expires,
			),
		}
		return assets, expires, postFetch, nil
	}
}
