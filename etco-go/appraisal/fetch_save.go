package appraisal

import (
	"fmt"
	"time"

	"github.com/WiggidyW/chanresult"
	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/purchasequeue"
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/shopassets"
)

func userMake(
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
		return appraisal, status, err
	} else if numActive >= build.PURCHASE_MAX_ACTIVE {
		status = MPS_MaxActiveLimit
		return appraisal, status, nil
	}

	// make the purchase
	status = MPS_Success
	err = remotedb.SetShopAppraisal(x, appraisal)
	return appraisal, status, err
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
