package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/cache/localcache"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

const (
	S_APPRAISAL_BUF_CAP          int           = 0
	S_APPRAISAL_LOCK_TTL         time.Duration = 30 * time.Second
	S_APPRAISAL_LOCK_MAX_BACKOFF time.Duration = 10 * time.Second
	S_APPRAISAL_EXPIRES_IN       time.Duration = 48 * time.Hour
)

func init() {
	keys.TypeStrShopAppraisal = localcache.RegisterType[ShopAppraisal](S_APPRAISAL_BUF_CAP)
}

type ShopAppraisal struct {
	// ignored during reading (used as doc id instead of field)
	// technically, if you're reading, you must already know it
	Code string `firestore:"-"`

	// ignored during writing (we use a nifty serverTimestamp firestore feature instead)
	Time time.Time `firestore:"time"`

	Items       []ShopItem `firestore:"items"`
	Price       float64    `firestore:"price"`
	TaxRate     float64    `firestore:"tax_rate,omitempty"`
	Tax         float64    `firestore:"tax,omitempty"`
	Version     string     `firestore:"version"`
	LocationId  int64      `firestore:"location_id"`
	CharacterId int32      `firestore:"character_id"`
}

func (s ShopAppraisal) GetCode() string {
	return s.Code
}

type ShopItem struct {
	TypeId       int32   `firestore:"type_id"`
	Quantity     int64   `firestore:"quantity"`
	PricePerUnit float64 `firestore:"price_per_unit"`
	Description  string  `firestore:"description"`
}

func GetShopAppraisal(
	ctx context.Context,
	code string,
) (
	rep *ShopAppraisal,
	expires *time.Time,
	err error,
) {
	return appraisalGet(
		ctx,
		client.readShopAppraisal,
		keys.TypeStrShopAppraisal,
		code,
		S_APPRAISAL_LOCK_TTL,
		S_APPRAISAL_LOCK_MAX_BACKOFF,
		S_APPRAISAL_EXPIRES_IN,
	)
}

func SetShopAppraisal(
	ctx context.Context,
	appraisal ShopAppraisal,
) (
	err error,
) {
	return appraisalSet(
		ctx,
		client.saveShopAppraisal,
		keys.TypeStrShopAppraisal,
		S_APPRAISAL_LOCK_TTL,
		S_APPRAISAL_LOCK_MAX_BACKOFF,
		S_APPRAISAL_EXPIRES_IN,
		appraisal,
		&[]prefetch.CacheAction{
			prefetch.ServerCacheDel(
				keys.TypeStrPurchaseQueue,
				keys.CacheKeyPurchaseQueue,
				S_APPRAISAL_LOCK_TTL,
				S_APPRAISAL_LOCK_MAX_BACKOFF,
			),
			prefetch.ServerCacheDel(
				keys.TypeStrUserData,
				keys.CacheKeyUserData(appraisal.CharacterId),
				S_APPRAISAL_LOCK_TTL,
				S_APPRAISAL_LOCK_MAX_BACKOFF,
			),
			prefetch.ServerCacheDel(
				keys.TypeStrUserShopAppraisalCodes,
				keys.CacheKeyUserShopAppraisalCodes(appraisal.CharacterId),
				S_APPRAISAL_LOCK_TTL,
				S_APPRAISAL_LOCK_MAX_BACKOFF,
			),
		},
	)
}
