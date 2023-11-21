package market

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

func GetShopPrice(
	x cache.Context,
	typeId int32,
	quantity int64,
	locationInfo staticdb.ShopLocationInfo,
) (
	price remotedb.ShopItem,
	expires time.Time,
	err error,
) {
	pricingInfo := locationInfo.GetTypePricingInfo(typeId)
	if pricingInfo == nil {
		price = newRejected(typeId, quantity)
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
		price = unpackPositivePrice(
			typeId,
			quantity,
			positivePrice,
			*pricingInfo,
		)
	}
	return price, expires, nil
}

func ProtoGetShopPrice(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	typeId int32,
	quantity int64,
	locationInfo staticdb.ShopLocationInfo,
) (
	price *proto.ShopItem,
	expires time.Time,
	err error,
) {
	var rPrice remotedb.ShopItem
	rPrice, expires, err = GetShopPrice(x, typeId, quantity, locationInfo)
	if err != nil {
		return nil, expires, err
	}
	return rPrice.ToProto(r), expires, nil
}

func NewRejectedShopItem(
	typeId int32,
	quantity int64,
) remotedb.ShopItem {
	return newRejected(typeId, quantity)
}

func newRejected(typeId int32, quantity int64) remotedb.ShopItem {
	return remotedb.ShopItem{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		Description:  Rejected(),
	}
}

func newRejectedNoOrders(
	typeId int32,
	quantity int64,
	mrktName string,
) remotedb.ShopItem {
	return remotedb.ShopItem{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		Description:  RejectedNoOrders(mrktName),
	}
}

func newAccepted(
	typeId int32,
	quantity int64,
	price float64,
	priceInfo staticdb.PricingInfo,
) remotedb.ShopItem {
	return remotedb.ShopItem{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: RoundedPrice(price),
		Description: Accepted(
			priceInfo.MarketName,
			priceInfo.Percentile,
			priceInfo.Modifier,
			priceInfo.IsBuy,
		),
	}
}

func unpackPositivePrice(
	typeId int32,
	quantity int64,
	positivePrice float64,
	priceInfo staticdb.PricingInfo,
) remotedb.ShopItem {
	accepted := positivePrice > 0.0
	if accepted {
		return newAccepted(typeId, quantity, positivePrice, priceInfo)
	} else {
		return newRejectedNoOrders(
			typeId,
			quantity,
			priceInfo.MarketName,
		)
	}
}
