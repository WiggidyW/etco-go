package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
)

func purchaseQueueCancel(
	x cache.Context,
	method func(context.Context) error,
	cacheLocks []cacheprefetch.ActionOrderedLocks,
	locationIds ...int64,
) (
	err error,
) {
	numLocIds := len(locationIds)
	locPurchaseQueueLocks := make([]cacheprefetch.ActionLock, 0, numLocIds)
	locUnreservedAssetLocks := make([]cacheprefetch.ActionLock, 0, numLocIds)
	for _, locationId := range locationIds {
		locPurchaseQueueLocks = append(
			locPurchaseQueueLocks,
			cacheprefetch.ServerLock(
				keys.CacheKeyLocationPurchaseQueue(locationId),
				keys.TypeStrLocationPurchaseQueue,
			),
		)
		locUnreservedAssetLocks = append(
			locUnreservedAssetLocks,
			cacheprefetch.ServerLock(
				keys.CacheKeyUnreservedShopAssets(locationId),
				keys.TypeStrUnreservedShopAssets,
			),
		)
	}
	_, _, err = fetch.FetchWithCache(
		x,
		purchaseQueueCancelFetchFunc(method),
		cacheprefetch.AntiCache[struct{}](append(
			cacheLocks,
			cacheprefetch.ActionOrderedLocks{
				Locks: locUnreservedAssetLocks,
				Child: &cacheprefetch.ActionOrderedLocks{
					Locks: locPurchaseQueueLocks,
					Child: &cacheprefetch.ActionOrderedLocks{
						Locks: []cacheprefetch.ActionLock{
							cacheprefetch.ServerLock(
								keys.CacheKeyPurchaseQueue,
								keys.TypeStrPurchaseQueue,
							),
						},
						Child: &cacheprefetch.ActionOrderedLocks{
							Locks: []cacheprefetch.ActionLock{
								cacheprefetch.ServerLock(
									keys.CacheKeyRawPurchaseQueue,
									keys.TypeStrRawPurchaseQueue,
								),
							},
							Child: nil,
						},
					},
				},
			},
		)),
	)
	return err
}

func purchaseQueueCancelFetchFunc(
	method func(context.Context) error,
) fetch.CachingFetch[struct{}] {
	return func(x cache.Context) (
		_ struct{},
		expires time.Time,
		_ *cachepostfetch.Params,
		err error,
	) {
		return struct{}{}, expires, nil, method(x.Ctx())
	}
}

func rawPurchaseQueueGet(x cache.Context) (
	rep RawPurchaseQueue,
	expires time.Time,
	err error,
) {
	return fetch.FetchWithCache(
		x,
		rawPurchaseQueueGetFetchFunc,
		cacheprefetch.StrongCache[RawPurchaseQueue](
			keys.CacheKeyRawPurchaseQueue,
			keys.TypeStrRawPurchaseQueue,
			nil,
			nil,
		),
	)
}

func rawPurchaseQueueGetFetchFunc(
	x cache.Context,
) (
	rep RawPurchaseQueue,
	expires time.Time,
	postFetch *cachepostfetch.Params,
	err error,
) {
	rep, err = readPurchaseQueue(x.Ctx())
	if err != nil {
		return nil, expires, nil, err
	}
	expires = time.Now().Add(FULL_PURCHASE_QUEUE_EXPIRES_IN)
	postFetch = &cachepostfetch.Params{
		Set: cachepostfetch.ServerSetOne[RawPurchaseQueue](
			keys.CacheKeyRawPurchaseQueue,
			keys.TypeStrRawPurchaseQueue,
			rep,
			expires,
		),
	}
	return rep, expires, postFetch, nil
}
