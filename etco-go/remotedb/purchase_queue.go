package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
	"github.com/WiggidyW/etco-go/remotedb/implrdb"
)

const (
	FULL_PURCHASE_QUEUE_EXPIRES_IN time.Duration = 24 * time.Hour
	FULL_PURCHASE_QUEUE_BUF_CAP    int           = 0
)

func init() {
	keys.TypeStrRawPurchaseQueue = cache.RegisterType[RawPurchaseQueue]("rawpurchasequeue", FULL_PURCHASE_QUEUE_BUF_CAP)
}

type RawPurchaseQueue = implrdb.RawPurchaseQueue

func GetRawPurchaseQueue(x cache.Context) (
	rep RawPurchaseQueue,
	expires time.Time,
	err error,
) {
	return rawPurchaseQueueGet(x)
}

func DelPurchases[C ICodeAndLocationId](
	x cache.Context,
	codes ...C,
) (
	err error,
) {
	locationIds := make([]int64, len(codes))
	for i, codeAndLocationId := range codes {
		locationIds[i] = codeAndLocationId.GetLocationId()
	}
	return purchaseQueueCancel(
		x,
		func(ctx context.Context) error {
			return delShopPurchases(ctx, codes...)
		},
		nil,
		locationIds...,
	)
}

func UserCancelPurchase(
	x cache.Context,
	characterId int32,
	code string,
	locationId int64,
) (
	err error,
) {
	return purchaseQueueCancel(
		x,
		func(ctx context.Context) error {
			return cancelShopPurchase(ctx, characterId, code, locationId)
		},
		[]cacheprefetch.ActionOrderedLocks{{
			Locks: []cacheprefetch.ActionLock{
				cacheprefetch.ServerLock(
					keys.CacheKeyUserCancelledPurchase(characterId),
					keys.TypeStrUserCancelledPurchase,
				),
			},
			Child: nil,
		}},
		locationId,
	)
}
