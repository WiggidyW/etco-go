package purchasequeue

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/items"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"
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

type LocationPurchaseQueue []string

func (lpq LocationPurchaseQueue) ToProto(
	locationInfo *proto.LocationInfo,
) *proto.PurchaseQueue {
	return &proto.PurchaseQueue{
		Codes:        lpq,
		LocationInfo: locationInfo,
	}
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
		rep = make(map[int32]int64, len(codeItems))
		items.AddToMap(rep, codeItems...)
		return rep, expires, err
	}

	// multi-fetch shop items
	x, cancel := x.WithCancel()
	defer cancel() // cancel if we return early
	chn := expirable.NewChanResult[[]remotedb.ShopItem](x.Ctx(), numCodes, 0)

	for _, code := range purchaseQueue {
		// Get the shop items in a goroutine
		go expirable.P2Transceive(
			chn,
			x, code,
			remotedb.GetShopAppraisalItems,
		)
	}

	// recv shop items and add them to the map
	var codeItems []remotedb.ShopItem
	for i := 0; i < numCodes; i++ {
		codeItems, _, err = chn.RecvExp()
		if i == 0 {
			rep = make(map[int32]int64, len(codeItems))
		}
		if err != nil {
			return nil, expires, err
		}
		items.AddToMap(rep, codeItems...)
	}

	return rep, expires, nil
}

func ProtoGetLocationPurchaseQueue(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	locationId int64,
) (
	rep *proto.PurchaseQueue,
	expires time.Time,
	err error,
) {
	// fetch proto location info
	x, cancel := x.WithCancel()
	defer cancel()
	locationInfoCOV := esi.ProtoGetLocationInfoCOV(x, r, locationId)

	// fetch location purchase queue
	var rQueue LocationPurchaseQueue
	rQueue, expires, err = GetLocationPurchaseQueue(x, locationId)
	if err != nil {
		return nil, expires, err
	}

	// recv proto location info
	var locationInfo *proto.LocationInfo
	locationInfo, expires, err = locationInfoCOV.RecvExpMin(expires)
	if err != nil {
		return nil, expires, err
	}

	return rQueue.ToProto(locationInfo), expires, nil
}

type PurchaseQueue map[int64][]string

func (pq PurchaseQueue) ToProto(
	locationInfos map[int64]*proto.LocationInfo,
) map[int64]*proto.PurchaseQueue {
	protoPQ := make(map[int64]*proto.PurchaseQueue, len(pq))
	for locationId, codes := range pq {
		protoPQ[locationId] =
			LocationPurchaseQueue(codes).ToProto(locationInfos[locationId])
	}
	return protoPQ
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

func InPurchaseQueue(
	x cache.Context,
	code string,
) (
	in bool,
	expires time.Time,
	err error,
) {
	var queue PurchaseQueue
	queue, expires, err = GetPurchaseQueue(x)
	if err != nil {
		return false, expires, err
	}
	for _, qCodes := range queue {
		for _, qCode := range qCodes {
			if qCode == code {
				return true, expires, nil
			}
		}
	}
	return false, expires, nil
}

func ProtoGetPurchaseQueue(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
) (
	rep map[int64]*proto.PurchaseQueue,
	expires time.Time,
	err error,
) {
	// fetch purchase queue
	var rQueue PurchaseQueue
	rQueue, expires, err = GetPurchaseQueue(x)
	if err != nil {
		return nil, expires, err
	}

	// fetch proto location info for each location in a goroutine
	x, cancel := x.WithCancel()
	defer cancel()
	chnInfo :=
		expirable.NewChanResult[*proto.LocationInfo](x.Ctx(), len(rQueue), 0)
	for locationId := range rQueue {
		go expirable.P3Transceive(
			chnInfo,
			x, r, locationId,
			esi.ProtoGetLocationInfo,
		)
	}

	// recv proto location info
	infoMap := make(map[int64]*proto.LocationInfo, len(rQueue))
	var locationInfo *proto.LocationInfo
	for i := 0; i < len(rQueue); i++ {
		locationInfo, expires, err = chnInfo.RecvExpMin(expires)
		if err != nil {
			return nil, expires, err
		} else if locationInfo != nil {
			infoMap[locationInfo.LocationId] = locationInfo
		}
	}

	return rQueue.ToProto(infoMap), expires, nil
}

func TransceiveGetPurchaseQueue(
	x cache.Context,
) (
	chnRecv expirable.ChanResult[PurchaseQueue],
) {
	return get(x)
}

func DelPurchases[C remotedb.ICodeAndLocationId](
	x cache.Context,
	codes ...C,
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

func ProtoUserCancelPurchase(
	x cache.Context,
	characterId int32,
	code string,
	locationId int64,
) (
	status proto.CancelPurchaseStatus,
	err error,
) {
	var rStatus CancelPurchaseStatus
	rStatus, err = userCancel(x, characterId, code, locationId)
	return rStatus.ToProto(), err
}
