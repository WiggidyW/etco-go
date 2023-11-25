package remotedb

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
	"github.com/WiggidyW/etco-go/logger"
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
	var rdbRep fsPurchaseQueue
	rdbRep, err = client.readPurchaseQueue(x.Ctx())
	if err != nil {
		return nil, expires, nil, err
	}
	expires = time.Now().Add(FULL_PURCHASE_QUEUE_EXPIRES_IN)
	rep = make(RawPurchaseQueue, len(rdbRep))
	for k, v := range rdbRep {
		// TODO: move this logic to the client, it should return map[int64][]string
		locationId, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			logger.Err(fmt.Sprintf(
				"Bad PurchaseQueue Key: [key: %s, err: %s]",
				k,
				err.Error(),
			))
			continue
		}
		interfaceCodes, ok := v.([]interface{})
		if !ok {
			logger.Err(fmt.Sprintf(
				"Bad PurchaseQueue Value: [key: %s, value: %v]",
				k,
				v,
			))
			continue
		}
		codes := make([]string, 0, len(interfaceCodes))
		for _, interfaceCode := range interfaceCodes {
			code, ok := interfaceCode.(string)
			if !ok {
				logger.Err(fmt.Sprintf(
					"Bad PurchaseQueue Code: [key: %s, value: %v]",
					k,
					interfaceCode,
				))
				continue
			}
			codes = append(codes, code)
		}
		rep[locationId] = codes
	}
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
