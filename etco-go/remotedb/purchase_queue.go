package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/cache/localcache"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

const (
	PURCHASE_QUEUE_BUF_CAP          int           = 0
	PURCHASE_QUEUE_LOCK_TTL         time.Duration = 1 * time.Minute
	PURCHASE_QUEUE_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	PURCHASE_QUEUE_EXPIRES_IN       time.Duration = 24 * time.Hour

	CANCEL_PURCHASE_LOCK_TTL         time.Duration = 1 * time.Minute
	CANCEL_PURCHASE_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute

	DEL_PURCHASES_LOCK_TTL         time.Duration = 1 * time.Minute
	DEL_PURCHASES_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
)

func init() {
	keys.TypeStrPurchaseQueue = localcache.RegisterType[[]string](PURCHASE_QUEUE_BUF_CAP)
}

type PurchaseQueue struct {
	PurchaseQueue []string `firestore:"shop_queue"`
}

func GetPurchaseQueue(ctx context.Context) (
	rep []string,
	expires *time.Time,
	err error,
) {
	return purchaseQueueGet(
		ctx,
		PURCHASE_QUEUE_LOCK_TTL,
		PURCHASE_QUEUE_LOCK_MAX_BACKOFF,
		PURCHASE_QUEUE_EXPIRES_IN,
	)
}

func DelPurchases(
	ctx context.Context,
	codes ...string,
) (
	err error,
) {
	return purchaseQueueCancel(
		ctx,
		func(context.Context) error {
			return client.delShopPurchases(ctx, codes...)
		},
		DEL_PURCHASES_LOCK_TTL,
		DEL_PURCHASES_LOCK_MAX_BACKOFF,
		nil,
	)
}

func UserCancelPurchase(
	ctx context.Context,
	characterId int32,
	code string,
) (
	err error,
) {
	cacheKeyUserData := keys.CacheKeyUserData(characterId)
	return purchaseQueueCancel(
		ctx,
		func(context.Context) error {
			return client.cancelShopPurchase(ctx, characterId, code)
		},
		CANCEL_PURCHASE_LOCK_TTL,
		CANCEL_PURCHASE_LOCK_MAX_BACKOFF,
		&[]prefetch.CacheAction{
			prefetch.ServerCacheDel(
				keys.TypeStrUserData,
				cacheKeyUserData,
				CANCEL_PURCHASE_LOCK_TTL,
				CANCEL_PURCHASE_LOCK_MAX_BACKOFF,
			),
			prefetch.ServerCacheDel(
				keys.TypeStrUserCancelledPurchase,
				cacheKeyUserData,
				CANCEL_PURCHASE_LOCK_TTL,
				CANCEL_PURCHASE_LOCK_MAX_BACKOFF,
			),
		},
	)
}
