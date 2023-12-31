package appraisal

import (
	"time"

	"github.com/WiggidyW/etco-go/appraisalcode"
	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/items"
)

type TerritoryId[UTERID any] interface {
	ToInt64() int64
	Unwrap() UTERID
}

type TerritoryInfo interface {
	GetTaxRate() float64
	GetFeePerM3() float64
}

type AppraisalItem interface {
	items.IBasicItem
	GetPricePerUnit() float64
	GetDescription() string
	GetFeePerUnit() float64
	GetChildrenLength() int
}

type priceTaxOperation uint8

const (
	PTAX_ADD priceTaxOperation = iota
	PTAX_SUB
	PTAX_NONE
)

type getTerritoryInfo[
	TERID any,
	TERINFO any,
] func(
	territoryId TERID,
) (
	territoryInfo *TERINFO,
)

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
	TERID any,
] func(
	rejected bool,
	code string,
	timeStamp time.Time,
	items []AITEM,
	version string,
	characterIdPtr *int32,
	territoryId TERID,
	price, tax, taxRate, fee, feePerM3 float64,
) (
	appraisal A,
)

func create[
	A any,
	AITEM AppraisalItem,
	BITEM items.IBasicItem,
	TERID TerritoryId[UTERID],
	UTERID any,
	TERINFO TerritoryInfo,
](
	x cache.Context,
	getTerritoryInfo getTerritoryInfo[UTERID, TERINFO],
	getPriceItem getPriceItem[AITEM, TERINFO],
	newAppraisal newAppraisal[A, AITEM, UTERID],
	includeCode *appraisalcode.CodeChar,
	priceTaxOperation priceTaxOperation,
	basicItems []BITEM,
	characterIdPtr *int32,
	territoryId TERID,
) (
	territoryInfo TERINFO,
	fee float64,
	appraisal A,
	expires time.Time,
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
	var price, tax, taxRate /* fee, */, feePerM3 float64
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
				territoryId.Unwrap(),
				price, tax, taxRate, fee, feePerM3,
			)
		}
	}()

	// get territory info and return rejected if it doesn't exist
	territoryInfoPtr := getTerritoryInfo(territoryId.Unwrap())
	if territoryInfoPtr == nil {
		rejected = true
		timeStamp = time.Now()
		expires = fetch.MAX_EXPIRES
		return territoryInfo, fee, appraisal, expires, nil
	}
	territoryInfo = *territoryInfoPtr
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
			return territoryInfo, fee, appraisal, expires, err
		} else {
			items = append(items, item)
		}
	}

	// finalize the appraisal
	timeStamp = time.Now()
	if priceTaxOperation == PTAX_ADD {
		tax = price * taxRate
		price += tax
	} else if priceTaxOperation == PTAX_SUB {
		tax = price * taxRate
		price -= tax
	} else /* if priceTaxOperation == PTAX_NONE */ {
		tax = 0.0
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
			territoryId.ToInt64(),
			price, tax, taxRate, fee, feePerM3,
		)
	}

	return territoryInfo, fee, appraisal, expires, nil
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
	expires = prevExpires
	price = prevPrice
	fee = prevFee

	item, expires, err = chn.RecvExpMin(expires)
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
