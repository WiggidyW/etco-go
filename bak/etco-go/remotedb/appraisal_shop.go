package remotedb

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

const (
	S_APPRAISAL_EXPIRES_IN time.Duration = 48 * time.Hour
	S_APPRAISAL_BUF_CAP    int           = 0
)

func init() {
	keys.TypeStrShopAppraisal = cache.RegisterType[ShopAppraisal]("shopappraisal", S_APPRAISAL_BUF_CAP)
}

type ShopAppraisal struct {
	Rejected bool `firestore:"rejected,omitempty"`

	// ignored during reading (used as doc id instead of field)
	// technically, if you're reading, you must already know it
	Code string `firestore:"-"`

	Time        time.Time  `firestore:"time"`
	Items       []ShopItem `firestore:"items"`
	Version     string     `firestore:"version"`
	CharacterId *int32     `firestore:"character_id"`
	LocationId  int64      `firestore:"location_id"`
	Price       float64    `firestore:"price"`
	Tax         float64    `firestore:"tax,omitempty"`
	TaxRate     float64    `firestore:"tax_rate,omitempty"`
}

func NewShopAppraisal(
	rejected bool,
	code string,
	timeStamp time.Time,
	items []ShopItem,
	version string,
	characterId *int32,
	locationId int64,
	price, tax, taxRate, _, _ float64,
) ShopAppraisal {
	return ShopAppraisal{
		Rejected:    rejected,
		Code:        code,
		Time:        timeStamp,
		Items:       items,
		Version:     version,
		CharacterId: characterId,
		LocationId:  locationId,
		Price:       price,
		Tax:         tax,
		TaxRate:     taxRate,
	}
}

func (s ShopAppraisal) GetCode() string { return s.Code }

type ShopItem struct {
	TypeId       int32   `firestore:"type_id"`
	Quantity     int64   `firestore:"quantity"`
	PricePerUnit float64 `firestore:"price_per_unit"`
	Description  string  `firestore:"description"`
}

func (si ShopItem) GetTypeId() int32         { return si.TypeId }
func (si ShopItem) GetQuantity() int64       { return si.Quantity }
func (si ShopItem) GetPricePerUnit() float64 { return si.PricePerUnit }
func (si ShopItem) GetDescription() string   { return si.Description }
func (si ShopItem) GetFeePerUnit() float64   { return 0.0 }
func (si ShopItem) GetChildrenLength() int   { return 0 }

func GetShopAppraisalItems(
	x cache.Context,
	code string,
) (
	rep []ShopItem,
	expires time.Time,
	err error,
) {
	var appraisal *ShopAppraisal
	appraisal, expires, err = GetShopAppraisal(x, code)
	if appraisal != nil {
		rep = appraisal.Items
	}
	return rep, expires, err
}

func GetShopAppraisal(
	x cache.Context,
	code string,
) (
	rep *ShopAppraisal,
	expires time.Time,
	err error,
) {
	return appraisalGet(
		x,
		client.readShopAppraisal,
		keys.TypeStrShopAppraisal,
		code,
		S_APPRAISAL_EXPIRES_IN,
	)
}

func SetShopAppraisal(
	x cache.Context,
	appraisal ShopAppraisal,
) (
	err error,
) {
	var cacheLocks []prefetch.CacheActionOrderedLocks
	if appraisal.CharacterId != nil {
		cacheLocks = []prefetch.CacheActionOrderedLocks{
			prefetch.CacheOrderedLocks(
				nil,
				prefetch.ServerCacheLock(
					keys.TypeStrUserShopAppraisalCodes,
					keys.CacheKeyUserShopAppraisalCodes(*appraisal.CharacterId),
				),
			),
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
					prefetch.ServerCacheLock(
						keys.TypeStrLocationPurchaseQueue,
						keys.CacheKeyLocationPurchaseQueue(appraisal.LocationId),
					),
				),
				prefetch.ServerCacheLock(
					keys.TypeStrUnreservedShopAssets,
					keys.CacheKeyUnreservedShopAssets(appraisal.LocationId),
				),
			),
		}
	}
	return appraisalSet(
		x,
		client.saveShopAppraisal,
		keys.TypeStrShopAppraisal,
		S_APPRAISAL_EXPIRES_IN,
		appraisal,
		cacheLocks,
	)
}
