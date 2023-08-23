package shop

import (
	"github.com/WiggidyW/weve-esi/client/appraisal"
	"github.com/WiggidyW/weve-esi/client/market"
	"github.com/WiggidyW/weve-esi/client/market/internal"
	"github.com/WiggidyW/weve-esi/staticdb"
)

type ShopPrice = appraisal.ShopItem

func newRejected(typeId int32, quantity int64) *ShopPrice {
	return &ShopPrice{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		Description:  market.Rejected(),
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
		Description:  market.RejectedNoOrders(mrktName),
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
		PricePerUnit: market.RoundedPrice(price),
		Description: market.Accepted(
			priceInfo.MrktName,
			priceInfo.Prctile,
			priceInfo.Modifier,
			priceInfo.IsBuy,
		),
	}
}

func unpackPositivePrice(
	typeId int32,
	quantity int64,
	positivePrice *internal.PositivePrice,
	priceInfo staticdb.PricingInfo,
) *ShopPrice {
	accepted, price := market.UnpackPositivePrice(positivePrice)
	if accepted {
		return newAccepted(typeId, quantity, price, priceInfo)
	} else {
		return newRejectedNoOrders(
			typeId,
			quantity,
			priceInfo.MrktName,
		)
	}
}
