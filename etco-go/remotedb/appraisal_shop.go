package remotedb

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
	"github.com/WiggidyW/etco-go/remotedb/implrdb"
)

const (
	S_APPRAISAL_EXPIRES_IN time.Duration = 48 * time.Hour
	S_APPRAISAL_BUF_CAP    int           = 0
)

func init() {
	keys.TypeStrShopAppraisal = cache.RegisterType[ShopAppraisal]("shopappraisal", S_APPRAISAL_BUF_CAP)
}

type ShopAppraisal = implrdb.ShopAppraisal
type ShopItem = implrdb.ShopItem

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
		readShopAppraisal,
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
	var cacheLocks []cacheprefetch.ActionOrderedLocks
	if appraisal.CharacterId != nil {
		cacheLocks = []cacheprefetch.ActionOrderedLocks{
			{
				Locks: []cacheprefetch.ActionLock{
					cacheprefetch.ServerLock(
						keys.CacheKeyUserShopAppraisalCodes(
							*appraisal.CharacterId,
						),
						keys.TypeStrUserShopAppraisalCodes,
					),
				},
				Child: nil,
			},
			{
				Locks: []cacheprefetch.ActionLock{
					cacheprefetch.ServerLock(
						keys.CacheKeyUnreservedShopAssets(appraisal.LocationId),
						keys.TypeStrUnreservedShopAssets,
					),
				},
				Child: &cacheprefetch.ActionOrderedLocks{
					Locks: []cacheprefetch.ActionLock{
						cacheprefetch.ServerLock(
							keys.CacheKeyLocationPurchaseQueue(
								appraisal.LocationId,
							),
							keys.TypeStrLocationPurchaseQueue,
						),
					},
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
		}
	}
	return appraisalSet(
		x,
		saveShopAppraisal,
		keys.TypeStrShopAppraisal,
		S_APPRAISAL_EXPIRES_IN,
		appraisal,
		cacheLocks,
	)
}
