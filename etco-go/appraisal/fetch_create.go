package appraisal

import (
	"time"

	"github.com/WiggidyW/etco-go/appraisalcode"
	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/items"
)

type taxOperation bool

const (
	TAX_ADD taxOperation = true
	TAX_SUB taxOperation = false
)

type getTerritoryInfo[TERID any, TERINFO any] func(territoryId TERID) *TERINFO

type getPriceItem[
	AITEM any,
	TERINFO any,
] func(
	x cache.Context,
	typeId int32,
	quantity int64,
	territoryInfo TERINFO,
) (
	price AITEM,
	expires time.Time,
	err error,
)

type newAppraisal[
	A any,
	AITEM any,
	TERID ~int64 | ~int32,
] func(
	rejected bool,
	code string,
	timeStamp time.Time,
	items []AITEM,
	version string,
	characterIdPtr *int32,
	territoryId TERID,
	price, tax, taxRate, fee, feePerM3 float64,
) A

func create[
	A any,
	AITEM AppraisalItem,
	BITEM items.IBasicItem,
	TERID ~int64 | ~int32,
	TERINFO TerritoryInfo,
](
	x cache.Context,
	getTerritoryInfo getTerritoryInfo[TERID, TERINFO],
	getPriceItem getPriceItem[AITEM, TERINFO],
	newAppraisal newAppraisal[A, AITEM, TERID],
	includeCode *appraisalcode.CodeChar,
	taxOperation taxOperation,
	basicItems []BITEM,
	characterIdPtr *int32,
	territoryId TERID,
) (
	appraisal A,
	expires time.Time,
	err error,
) {
	return fetch.HandleFetch(
		x,
		nil,
		createFetchFunc(
			getTerritoryInfo,
			getPriceItem,
			newAppraisal,
			includeCode,
			taxOperation,
			basicItems,
			characterIdPtr,
			territoryId,
		),
		nil,
	)
}

func createFetchFunc[
	A any,
	AITEM AppraisalItem,
	BITEM items.IBasicItem,
	TERID ~int64 | ~int32,
	TERINFO TerritoryInfo,
](
	getTerritoryInfo getTerritoryInfo[TERID, TERINFO],
	getPriceItem getPriceItem[AITEM, TERINFO],
	newAppraisal newAppraisal[A, AITEM, TERID],
	includeCode *appraisalcode.CodeChar,
	taxOperation taxOperation,
	basicItems []BITEM,
	characterIdPtr *int32,
	territoryId TERID,
) fetch.Fetch[A] {
	return func(x cache.Context) (
		appraisal A,
		expires time.Time,
		_ *postfetch.Params,
		err error,
	) {
		// // initialize our variables and defer a 'newAppraisal' function call
		// // this does look weird, but it's really convenient
		var rejected bool
		var code string
		var timeStamp time.Time
		var items []AITEM
		var version string = build.DATA_VERSION
		// var characterIdPtr *int32
		var price, tax, taxRate, fee, feePerM3 float64
		// var territoryId TERID
		defer func() {
			if err == nil {
				appraisal = newAppraisal(
					rejected,
					code,
					timeStamp,
					items,
					version,
					characterIdPtr,
					territoryId,
					price, tax, taxRate, fee, feePerM3,
				)
			}
		}()

		// get territory info and return rejected if it doesn't exist
		territoryInfoPtr := getTerritoryInfo(territoryId)
		if territoryInfoPtr == nil {
			rejected = true
			timeStamp = time.Now()
			expires = fetch.MAX_EXPIRES
			return appraisal, expires, nil, nil
		}
		territoryInfo := *territoryInfoPtr
		taxRate = territoryInfo.GetTaxRate()
		feePerM3 = territoryInfo.GetFeePerM3()

		// fetch the price items in parallel
		x, cancel := x.WithCancel()
		defer cancel()
		chn := expirable.NewChanResult[AITEM](x.Ctx(), len(basicItems), 0)
		for _, basicItem := range basicItems {
			go transceiveGetPriceItem(
				x,
				getPriceItem,
				basicItem,
				territoryInfo,
				chn,
			)
		}

		// collect the price items
		expires = fetch.MAX_EXPIRES
		items = make([]AITEM, 0, len(basicItems))
		var item AITEM
		for i := 0; i < len(basicItems); i++ {
			item, expires, price, fee, err =
				handleRecvPrice(chn, expires, price, fee)
			if err != nil {
				return appraisal, expires, nil, err
			} else {
				items = append(items, item)
			}
		}

		// finalize the appraisal
		timeStamp = time.Now()
		tax = price * taxRate
		if taxOperation == TAX_ADD {
			price += tax
		} else /* if taxOperation == TAX_SUB */ {
			price -= tax
		}
		if price <= 0.0 {
			rejected = true
		} else if includeCode != nil { // never hash if rejected
			code = hashAppraisal(
				*includeCode,
				timeStamp,
				len(items),
				version,
				characterIdPtr,
				territoryId,
				price, tax, taxRate, fee, feePerM3,
			)
		}

		return appraisal, expires, nil, nil
	}
}

func handleRecvPrice[AITEM AppraisalItem](
	chn expirable.ChanResult[AITEM],
	prevExpires time.Time,
	prevPrice float64,
	prevFee float64,
) (
	item AITEM,
	expires time.Time,
	price float64,
	fee float64,
	err error,
) {
	item, expires, err = chn.RecvExpMin(prevExpires)
	if err != nil {
		return item, expires, price, fee, err
	}

	itemFeePerUnit := item.GetFeePerUnit()
	itemSumPerUnit := item.GetPricePerUnit() - itemFeePerUnit
	if itemSumPerUnit > 0.0 {
		itemQuantityF64 := float64(item.GetQuantity())
		price += itemSumPerUnit * itemQuantityF64
		fee += itemFeePerUnit * itemQuantityF64
	}

	return item, expires, price, fee, nil
}

func transceiveGetPriceItem[AITEM any, BITEM items.IBasicItem, TERINFO any](
	x cache.Context,
	getPriceItem getPriceItem[AITEM, TERINFO],
	basicItem BITEM,
	territoryInfo TERINFO,
	chn expirable.ChanResult[AITEM],
) error {
	return chn.SendExp(getPriceItem(
		x,
		basicItem.GetTypeId(),
		basicItem.GetQuantity(),
		territoryInfo,
	))
}
