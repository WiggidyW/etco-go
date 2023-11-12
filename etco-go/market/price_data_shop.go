package market

import (
	rdb "github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

func NewRejectedShopItem(
	typeId int32,
	quantity int64,
) *rdb.ShopItem {
	return newRejected(typeId, quantity)
}

type ShopPrice = rdb.ShopItem

func newRejected(typeId int32, quantity int64) *ShopPrice {
	return &ShopPrice{
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
) *ShopPrice {
	return &ShopPrice{
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
) *ShopPrice {
	return &ShopPrice{
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
) *ShopPrice {
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
