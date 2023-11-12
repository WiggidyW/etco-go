package remotedb

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
	"github.com/WiggidyW/etco-go/logger"
)

func purchaseQueueCancel(
	x cache.Context,
	method func(context.Context) error,
	cacheLocks []prefetch.CacheActionOrderedLocks,
	locationIds ...int64,
) (
	err error,
) {
	numLocIds := len(locationIds)
	locPurchaseQueueLocks := make([]prefetch.CacheActionLock, 0, numLocIds)
	locUnreservedAssetLocks := make([]prefetch.CacheActionLock, 0, numLocIds)
	for _, locationId := range locationIds {
		locPurchaseQueueLocks = append(
			locPurchaseQueueLocks,
			prefetch.ServerCacheLock(
				keys.CacheKeyLocationPurchaseQueue(locationId),
				keys.TypeStrLocationPurchaseQueue,
			),
		)
		locUnreservedAssetLocks = append(
			locUnreservedAssetLocks,
			prefetch.ServerCacheLock(
				keys.CacheKeyUnreservedShopAssets(locationId),
				keys.TypeStrUnreservedShopAssets,
			),
		)
	}
	_, _, err = fetch.HandleFetch(
		x,
		&prefetch.Params[struct{}]{
			CacheParams: &prefetch.CacheParams[struct{}]{
				Lock: append(
					cacheLocks,
					prefetch.CacheOrderedLocks(
						prefetch.CacheOrderedLocksPtr(
							prefetch.CacheOrderedLocksPtr(
								prefetch.CacheOrderedLocksPtr(
									nil,
									prefetch.ServerCacheLock(
										keys.CacheKeyRawPurchaseQueue,
										keys.TypeStrRawPurchaseQueue,
									),
								),
								prefetch.ServerCacheLock(
									keys.CacheKeyPurchaseQueue,
									keys.TypeStrPurchaseQueue,
								),
							),
							locPurchaseQueueLocks...,
						),
						locUnreservedAssetLocks...,
					),
				),
			},
		},
		purchaseQueueCancelFetchFunc(method),
		nil,
	)
	return err
}

func purchaseQueueCancelFetchFunc(
	method func(context.Context) error,
) fetch.Fetch[struct{}] {
	return func(x cache.Context) (
		_ *struct{},
		expires time.Time,
		_ *postfetch.Params,
		err error,
	) {
		return nil, expires, nil, method(x.Ctx())
	}
}

func rawPurchaseQueueGet(x cache.Context) (
	rep RawPurchaseQueue,
	expires time.Time,
	err error,
) {
	return fetch.HandleFetchVal(
		x,
		&prefetch.Params[RawPurchaseQueue]{
			CacheParams: &prefetch.CacheParams[RawPurchaseQueue]{
				Get: prefetch.ServerCacheGet[RawPurchaseQueue](
					keys.CacheKeyRawPurchaseQueue,
					keys.TypeStrRawPurchaseQueue,
					true,
					nil,
				),
			},
		},
		rawPurchaseQueueGetFetchFunc,
		nil,
	)
}

func rawPurchaseQueueGetFetchFunc(
	x cache.Context,
) (
	repPtr *RawPurchaseQueue,
	expires time.Time,
	postFetch *postfetch.Params,
	err error,
) {
	var rdbRep fsPurchaseQueue
	rdbRep, err = client.readPurchaseQueue(x.Ctx())
	if err != nil {
		return nil, expires, nil, err
	}
	expires = time.Now().Add(FULL_PURCHASE_QUEUE_EXPIRES_IN)
	rep := make(RawPurchaseQueue, len(rdbRep))
	for k, v := range rdbRep {
		locationId, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			logger.Err(fmt.Sprintf(
				"Bad PurchaseQueue Key: [key: %s, err: %s]",
				k,
				err.Error(),
			))
			continue
		}
		codes, ok := v.([]string)
		if !ok {
			logger.Err(fmt.Sprintf(
				"Bad PurchaseQueue Value: [key: %s, value: %v]",
				k,
				v,
			))
			continue
		}
		rep[locationId] = codes
	}
	postFetch = &postfetch.Params{
		CacheParams: &postfetch.CacheParams{
			Set: postfetch.ServerCacheSetOne(
				keys.CacheKeyRawPurchaseQueue,
				keys.TypeStrRawPurchaseQueue,
				&rep,
				expires,
			),
		},
	}
	return &rep, expires, postFetch, nil
}
