package market

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

func GetHaulPrice(
	x cache.Context,
	typeId int32,
	quantity int64,
	routeInfo staticdb.HaulRouteInfo,
	fallbackPricePerUnit map[int32]float64,
) (
	price remotedb.HaulItem,
	expires time.Time,
	err error,
) {
	pricingInfo := routeInfo.GetTypePricingInfo(typeId)
	if pricingInfo == nil {
		price = newRejectedHaulItem(typeId, quantity)
		expires = fetch.MAX_EXPIRES
	} else {
		var positivePrice float64
		positivePrice, expires, err = GetPercentilePrice(
			x,
			typeId,
			*pricingInfo,
		)
		if err != nil {
			return price, expires, err
		}
		price = unpackPositiveHaulPrice(
			typeId,
			quantity,
			positivePrice,
			*pricingInfo,
			routeInfo.M3Fee,
			fallbackPricePerUnit,
		)
	}
	return price, expires, nil
}

func NewRejectedHaulItem(typeId int32, quantity int64) remotedb.HaulItem {
	return newRejectedHaulItem(typeId, quantity)
}

func newRejectedHaulItem(typeId int32, quantity int64) remotedb.HaulItem {
	return remotedb.HaulItem{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		FeePerUnit:   0.0,
		Description:  Rejected(),
	}
}

func newRejectedHaulItemNoOrders(
	typeId int32,
	quantity int64,
	marketName string,
) remotedb.HaulItem {
	return remotedb.HaulItem{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		FeePerUnit:   0.0,
		Description:  RejectedNoOrders(marketName),
	}
}

func newAcceptedHaulItem(
	typeId int32,
	quantity int64,
	price float64,
	fee float64,
	priceInfo staticdb.PricingInfo,
) remotedb.HaulItem {
	return remotedb.HaulItem{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: RoundedPrice(price),
		FeePerUnit:   fee,
		Description: Accepted(
			priceInfo.MarketName,
			priceInfo.Percentile,
			priceInfo.Modifier,
			priceInfo.IsBuy,
		),
	}
}

func newAcceptedFallbackHaulItem(
	typeId int32,
	quantity int64,
	price float64,
	fee float64,
) remotedb.HaulItem {
	return remotedb.HaulItem{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: RoundedPrice(price),
		FeePerUnit:   fee,
		Description:  AcceptedFallback(),
	}
}

func unpackPositiveHaulPrice(
	typeId int32,
	quantity int64,
	positivePrice float64,
	priceInfo staticdb.PricingInfo,
	feePerM3 float64,
	fallbackPricePerUnit map[int32]float64,
) remotedb.HaulItem {
	var fallback bool = false
	accepted := positivePrice > 0.0
	if !accepted {
		positivePrice = fallbackPricePerUnit[typeId]
		if positivePrice > 0.0 {
			fallback = true
		} else {
			return newRejectedHaulItemNoOrders(
				typeId,
				quantity,
				priceInfo.MarketName,
			)
		}
	}
	fee := calculateTypeFee(
		typeId,
		nil,
		feePerM3,
	)
	if fallback {
		return newAcceptedFallbackHaulItem(
			typeId,
			quantity,
			positivePrice,
			fee,
		)
	} else {
		return newAcceptedHaulItem(
			typeId,
			quantity,
			positivePrice,
			fee,
			priceInfo,
		)
	}
}
