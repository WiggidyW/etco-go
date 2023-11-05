package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/cache/localcache"
)

const (
	PURCHASE_QUEUE_BUF_CAP          int           = 0
	PURCHASE_QUEUE_LOCK_TTL         time.Duration = 1 * time.Minute
	PURCHASE_QUEUE_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	PURCHASE_QUEUE_EXPIRES_IN       time.Duration = 24 * time.Hour
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
