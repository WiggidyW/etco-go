package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/cache/localcache"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

const (
	B_APPRAISAL_BUF_CAP          int           = 0
	B_APPRAISAL_LOCK_TTL         time.Duration = 30 * time.Second
	B_APPRAISAL_LOCK_MAX_BACKOFF time.Duration = 10 * time.Second
	B_APPRAISAL_EXPIRES_IN       time.Duration = 48 * time.Hour
)

func init() {
	keys.TypeStrBuybackAppraisal = localcache.RegisterType[BuybackAppraisal](B_APPRAISAL_BUF_CAP)
}

type BuybackAppraisal struct {
	// ignored during reading (used as doc id instead of field)
	// technically, if you're reading, you must already know it
	Code string `firestore:"-"`

	// ignored during writing (we use a nifty serverTimestamp firestore feature instead)
	Time time.Time `firestore:"time"`

	Items       []BuybackParentItem `firestore:"items"`
	FeePerM3    float64             `firestore:"fee_per_m3,omitempty"`
	Fee         float64             `firestore:"fee,omitempty"`
	TaxRate     float64             `firestore:"tax_rate,omitempty"`
	Tax         float64             `firestore:"tax,omitempty"`
	Price       float64             `firestore:"price"`
	Version     string              `firestore:"version"`
	SystemId    int32               `firestore:"system_id"`
	CharacterId *int32              `firestore:"character_id"`
}

func (b BuybackAppraisal) GetCode() string {
	return b.Code
}

type BuybackParentItem struct {
	TypeId       int32              `firestore:"type_id"`
	Quantity     int64              `firestore:"quantity"`
	PricePerUnit float64            `firestore:"price_per_unit"`
	FeePerUnit   float64            `firestore:"fee,omitempty"`
	Description  string             `firestore:"description"`
	Children     []BuybackChildItem `firestore:"children"`
}

type BuybackChildItem struct {
	TypeId            int32   `firestore:"type_id"`
	QuantityPerParent float64 `firestore:"quantity_per_parent"`
	PricePerUnit      float64 `firestore:"price_per_unit"`
	Description       string  `firestore:"description"`
}

func GetBuybackAppraisal(
	ctx context.Context,
	code string,
) (
	rep *BuybackAppraisal,
	expires *time.Time,
	err error,
) {
	return appraisalGet(
		ctx,
		client.readBuybackAppraisal,
		keys.TypeStrBuybackAppraisal,
		code,
		B_APPRAISAL_LOCK_TTL,
		B_APPRAISAL_LOCK_MAX_BACKOFF,
		B_APPRAISAL_EXPIRES_IN,
	)
}

func SetBuybackAppraisal(
	ctx context.Context,
	appraisal BuybackAppraisal,
) (
	err error,
) {
	var cacheDels *[]prefetch.CacheAction
	if appraisal.CharacterId != nil {
		cacheDels = &[]prefetch.CacheAction{
			prefetch.ServerCacheDel(
				keys.TypeStrUserData,
				keys.CacheKeyUserData(*appraisal.CharacterId),
				B_APPRAISAL_LOCK_TTL,
				B_APPRAISAL_LOCK_MAX_BACKOFF,
			),
			prefetch.ServerCacheDel(
				keys.TypeStrUserBuybackAppraisalCodes,
				keys.CacheKeyUserBuybackAppraisalCodes(*appraisal.CharacterId),
				B_APPRAISAL_LOCK_TTL,
				B_APPRAISAL_LOCK_MAX_BACKOFF,
			),
		}
	}
	return appraisalSet(
		ctx,
		client.saveBuybackAppraisal,
		keys.TypeStrBuybackAppraisal,
		B_APPRAISAL_LOCK_TTL,
		B_APPRAISAL_LOCK_MAX_BACKOFF,
		B_APPRAISAL_EXPIRES_IN,
		appraisal,
		cacheDels,
	)
}
