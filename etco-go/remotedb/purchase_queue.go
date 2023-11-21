package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

const (
	FULL_PURCHASE_QUEUE_EXPIRES_IN time.Duration = 24 * time.Hour
	FULL_PURCHASE_QUEUE_BUF_CAP    int           = 0
)

func init() {
	keys.TypeStrRawPurchaseQueue = cache.RegisterType[RawPurchaseQueue]("rawpurchasequeue", FULL_PURCHASE_QUEUE_BUF_CAP)
}

type fsPurchaseQueue = map[string]interface{}

type RawPurchaseQueue = map[int64][]string

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
			return delShopPurchases(client, ctx, codes...)
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
			return client.cancelShopPurchase(ctx, characterId, code, locationId)
		},
		prefetch.ServerCacheOrderedLocksOne(
			keys.CacheKeyUserCancelledPurchase(characterId),
			keys.TypeStrUserCancelledPurchase,
		),
		locationId,
	)
}
