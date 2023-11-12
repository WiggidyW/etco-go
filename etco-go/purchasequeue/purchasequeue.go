package purchasequeue

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/items"
	"github.com/WiggidyW/etco-go/remotedb"
)

const (
	PURCHASE_QUEUE_BUF_CAP          int = 0
	LOCATION_PURCHASE_QUEUE_BUF_CAP int = 0
)

func init() {
	keys.TypeStrPurchaseQueue = cache.RegisterType[PurchaseQueue]("purchasequeue", PURCHASE_QUEUE_BUF_CAP)
	keys.TypeStrLocationPurchaseQueue = cache.RegisterType[LocationPurchaseQueue]("locationpurchasequeue", LOCATION_PURCHASE_QUEUE_BUF_CAP)
}

type LocationPurchaseQueue = []string
type PurchaseQueue = map[int64][]string

func GetLocationPurchaseQueueItems(
	x cache.Context,
	locationId int64,
) (
	rep map[int32]int64, // map[typeId]quantity
	expires time.Time,
	err error,
) {
	// !IMPORTANT: we will ignore appraisal expire times as they are arbitrary
	// !IMPORTANT: if that ever changes, we will need to update this

	var purchaseQueue LocationPurchaseQueue
	purchaseQueue, expires, err = locationGet(x, locationId)

	if err != nil || purchaseQueue == nil { // no fetches needed
		return nil, expires, err
	}

	numCodes := len(purchaseQueue)

	if numCodes == 1 { // no concurrency needed
		codeItems, _, err := remotedb.GetShopAppraisalItems(x, purchaseQueue[0])
		items.AddToMap(rep, codeItems...)
		return rep, expires, err
	}

	// multi-fetch shop items
	x, cancel := x.WithCancel()
	defer cancel() // cancel if we return early
	chn := expirable.NewChanResult[[]remotedb.ShopItem](x.Ctx(), numCodes, 0)

	for _, code := range purchaseQueue {
		// Get the shop items in a goroutine
		go expirable.Param2Transceive(
			chn,
			x, code,
			remotedb.GetShopAppraisalItems,
		)
	}

	// recv shop items and add them to the map
	var codeItems []remotedb.ShopItem
	for i := 0; i < numCodes; i++ {
		codeItems, _, err = chn.RecvExp()
		if err != nil {
			return nil, expires, err
		}
		items.AddToMap(rep, codeItems...)
	}

	return rep, expires, nil
}

func GetLocationPurchaseQueue(
	x cache.Context,
	locationId int64,
) (
	queue LocationPurchaseQueue,
	expires time.Time,
	err error,
) {
	return locationGet(x, locationId)
}

func GetPurchaseQueue(
	x cache.Context,
) (
	queue PurchaseQueue,
	expires time.Time,
	err error,
) {
	return TransceiveGetPurchaseQueue(x).RecvExp()
}

func TransceiveGetPurchaseQueue(
	x cache.Context,
) (
	chnRecv expirable.ChanResult[PurchaseQueue],
) {
	return get(x)
}

func DelPurchases(
	x cache.Context,
	codes ...remotedb.CodeAndLocationId,
) (
	err error,
) {
	return remotedb.DelPurchases(x, codes...)
}

func UserCancelPurchase(
	x cache.Context,
	characterId int32,
	code string,
	locationId int64,
) (
	status CancelPurchaseStatus,
	err error,
) {
	return userCancel(x, characterId, code, locationId)
}
