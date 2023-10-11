package market

import (
	"github.com/WiggidyW/etco-go/client/market/marketprice"
	rdb "github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

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
	positivePrice *marketprice.PositivePrice,
	priceInfo staticdb.PricingInfo,
) *ShopPrice {
	accepted, price := UnpackPositivePrice(positivePrice)
	if accepted {
		return newAccepted(typeId, quantity, price, priceInfo)
	} else {
		return newRejectedNoOrders(
			typeId,
			quantity,
			priceInfo.MarketName,
		)
	}
}
