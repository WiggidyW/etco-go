package appraisal

import (
	"context"
	"fmt"
	"time"

	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go/appraisalcode"
	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/items"
	"github.com/WiggidyW/etco-go/logger"
	"github.com/WiggidyW/etco-go/market"
	"github.com/WiggidyW/etco-go/notifier"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
	"github.com/WiggidyW/etco-go/protoregistry"
	"github.com/WiggidyW/etco-go/purchasequeue"
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/shopassets"
	"github.com/WiggidyW/etco-go/staticdb"
)

func GetShopAppraisal(
	x cache.Context,
	code string,
) (
	appraisal *remotedb.ShopAppraisal,
	expires time.Time,
	err error,
) {
	return remotedb.GetShopAppraisal(x, code)
}

func ProtoGetShopAppraisal(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
	include_items bool,
) (
	appraisal *proto.ShopAppraisal,
	expires time.Time,
	err error,
) {
	var rAppraisal *remotedb.ShopAppraisal
	rAppraisal, expires, err = GetShopAppraisal(x, code)
	if err != nil {
		return nil, expires, err
	} else if rAppraisal == nil {
		return nil, expires, protoerr.MsgNew(
			protoerr.NOT_FOUND,
			"Shop Appraisal not found",
		)
	} else if !include_items {
		rAppraisal.Items = nil
	}

	var locationInfo *proto.LocationInfo
	locationInfo, expires, err = esi.
		ProtoGetLocationInfoCOV(x, r, rAppraisal.LocationId).
		RecvExpMin(expires)
	if err != nil {
		return nil, expires, err
	}

	return rAppraisal.ToProto(r, locationInfo), expires, nil
}

func GetShopAppraisalCharacterId(
	x cache.Context,
	code string,
) (
	characterId *int32,
	expires time.Time,
	err error,
) {
	var appraisal *remotedb.ShopAppraisal
	appraisal, expires, err = GetShopAppraisal(x, code)
	if err == nil && appraisal != nil {
		characterId = appraisal.CharacterId
	}
	return characterId, expires, err
}

func CreateShopAppraisal[BITEM items.IBasicItem](
	x cache.Context,
	items []BITEM,
	characterId *int32,
	locationId int64,
	includeCode bool,
) (
	appraisal remotedb.ShopAppraisal,
	expires time.Time,
	err error,
) {
	var codeChar *appraisalcode.CodeChar
	if includeCode {
		codeChar = &appraisalcode.BUYBACK_CHAR
	}
	return create(
		x,
		staticdb.GetShopLocationInfo,
		market.GetShopPrice,
		remotedb.NewShopAppraisal,
		codeChar,
		TAX_ADD,
		items,
		characterId,
		locationId,
	)
}

func ProtoCreateShopAppraisal[BITEM items.IBasicItem](
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	items []BITEM,
	characterId *int32,
	locationId int64,
	save bool,
) (
	appraisal *proto.ShopAppraisal,
	status proto.MakePurchaseStatus,
	expires time.Time,
	err error,
) {
	x, cancel := x.WithCancel()
	defer cancel()
	locationInfoCOV := esi.ProtoGetLocationInfoCOV(x, r, locationId)

	var rAppraisal remotedb.ShopAppraisal
	var rStatus MakePurchaseStatus
	if save {
		rAppraisal, rStatus, expires, err = CreateSaveShopAppraisal(
			x,
			items,
			characterId,
			locationId,
		)
	} else {
		rStatus = MPS_None
		rAppraisal, expires, err = CreateShopAppraisal(
			x,
			items,
			characterId,
			locationId,
			false,
		)
	}
	if err != nil {
		return nil, status, expires, err
	} else {
		var locationInfo *proto.LocationInfo
		locationInfo, expires, err = locationInfoCOV.RecvExpMin(expires)
		if err != nil {
			return nil, status, expires, err
		}
		appraisal = rAppraisal.ToProto(r, locationInfo)
		status = rStatus.ToProto()
		return appraisal, status, expires, nil
	}
}

func SaveShopAppraisal(
	x cache.Context,
	appraisalIn remotedb.ShopAppraisal,
) (
	appraisal remotedb.ShopAppraisal,
	status MakePurchaseStatus,
	err error,
) {
	appraisal = appraisalIn
	if appraisal.CharacterId == nil {
		status = MPS_None
		err = remotedb.SetShopAppraisal(x, appraisal)
		return appraisal, status, err
	} else {
		return saveShopAppraisalAsPurchase(x, appraisal)
	}
}

func CreateSaveShopAppraisal[BITEM items.IBasicItem](
	x cache.Context,
	items []BITEM,
	characterId *int32,
	locationId int64,
) (
	appraisal remotedb.ShopAppraisal,
	status MakePurchaseStatus,
	expires time.Time,
	err error,
) {
	appraisal, expires, err = CreateShopAppraisal(
		x,
		items,
		characterId,
		locationId,
		true,
	)
	if err == nil && !appraisal.Rejected && appraisal.Price > 0.0 {
		appraisal, status, err = SaveShopAppraisal(x, appraisal)
	}
	return appraisal, status, expires, err
}

func saveShopAppraisalAsPurchase(
	x cache.Context,
	appraisalIn remotedb.ShopAppraisal,
) (
	appraisal remotedb.ShopAppraisal,
	status MakePurchaseStatus,
	err error,
) {
	appraisal = appraisalIn

	x, cancel := x.WithCancel()
	defer cancel()

	// fetch unreserved location assets in a goroutine
	chnAssets := expirable.NewChanResult[map[int32]int64](x.Ctx(), 1, 0)
	go expirable.P2Transceive(
		chnAssets,
		x, appraisal.LocationId,
		shopassets.GetUnreservedShopAssets,
	)

	// fetch num active purchases in a goroutine
	chnActive := chanresult.NewChanResult[int](x.Ctx(), 1, 0)
	go userActivePurchases(x, *appraisal.CharacterId, chnActive)

	// fetch user made purchase time and check cooldown
	var userMade *time.Time
	userMade, _, err = remotedb.GetUserMadePurchase(x, *appraisal.CharacterId)
	if err != nil {
		status = MPS_None
		return appraisal, status, err
	} else if userMade != nil && time.Now().Before(
		(*userMade).Add(build.MAKE_PURCHASE_COOLDOWN),
	) {
		status = MPS_CooldownLimit
		return appraisal, status, nil
	}

	// recv unreserved location assets and check if items are rejected or unavailable
	var assets map[int32]int64
	assets, _, err = chnAssets.RecvExp()
	if err != nil {
		status = MPS_None
		return appraisal, status, err
	}
	var ok bool
	status, ok = checkRejectedOrUnavailable(x, &appraisal, assets)
	if !ok {
		return appraisal, status, nil
	}

	// recv num active purchases and check if user has too many active purchases
	var numActive int
	numActive, err = chnActive.Recv()
	if err != nil {
		status = MPS_None
		return appraisal, status, err
	} else if numActive >= build.PURCHASE_MAX_ACTIVE {
		status = MPS_MaxActiveLimit
		return appraisal, status, nil
	}

	// make the purchase
	err = remotedb.SetShopAppraisal(x, appraisal)
	if err != nil {
		status = MPS_None
		return appraisal, status, err
	} else if build.PURCHASE_NOTIFICATIONS {
		go func() {
			logger.MaybeErr(notifier.PurchasesSend(
				context.Background(),
				appraisal.Code,
			))
		}()
	}

	status = MPS_Success
	return appraisal, status, nil
}

func checkRejectedOrUnavailable(
	x cache.Context,
	appraisal *remotedb.ShopAppraisal,
	assets map[int32]int64,
) (
	status MakePurchaseStatus,
	ok bool,
) {
	var rejected bool = false
	var unavailable bool = false

	for i := 0; i < len(appraisal.Items); i++ {
		item := &appraisal.Items[i]

		if item.PricePerUnit <= 0.0 {
			rejected = true
		}

		available := assets[item.TypeId]
		if available < item.Quantity {
			unavailable = true
			item.PricePerUnit = 0.0
			item.Description = fmt.Sprintf(
				"Rejected - %d are available for purchase",
				available,
			)
		}
	}

	if rejected && unavailable {
		ok = false
		status = MPS_ItemsRejectedAndUnavailable
	} else if rejected {
		ok = false
		status = MPS_ItemsRejected
	} else if unavailable {
		ok = false
		status = MPS_ItemsUnavailable
	} else {
		ok = true
		status = MPS_Success
	}

	return status, ok
}

func userActivePurchases(
	x cache.Context,
	characterId int32,
	chn chanresult.ChanResult[int],
) (
	ctxErr error,
) {
	var numActive int = 0
	var err error

	// fetch location purchase queue in a goroutine
	chnQueue :=
		expirable.NewChanResult[purchasequeue.PurchaseQueue](x.Ctx(), 1, 0)
	go expirable.P1Transceive(
		chnQueue,
		x,
		purchasequeue.GetPurchaseQueue,
	)

	// fetch user appraisal codes
	var userCodes []string
	userCodes, _, err = remotedb.GetUserShopAppraisalCodes(x, characterId)
	if err != nil {
		return chn.SendErr(err)
	}

	// convert them to a map
	userCodesSet := make(map[string]struct{}, len(userCodes))
	for _, code := range userCodes {
		userCodesSet[code] = struct{}{}
	}

	// recv location purchase queue
	var queue map[int64][]string
	queue, _, err = chnQueue.RecvExp()
	if err != nil {
		return chn.SendErr(err)
	}

	// count active purchases
	for _, codes := range queue {
		for _, code := range codes {
			if _, ok := userCodesSet[code]; ok {
				numActive++
			}
		}
	}

	return chn.SendOk(numActive)
}
