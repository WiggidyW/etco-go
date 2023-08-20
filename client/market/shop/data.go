package shop

import (
	"github.com/WiggidyW/weve-esi/client/market"
	"github.com/WiggidyW/weve-esi/client/market/internal"
	"github.com/WiggidyW/weve-esi/staticdb"
)

type ShopPrice = market.MarketPrice

func newRejected() *ShopPrice {
	return &ShopPrice{
		Price: 0.0,
		Desc:  market.Rejected(),
	}
}

func newRejectedNoOrders(mrktName string) *ShopPrice {
	return &ShopPrice{
		Price: 0.0,
		Desc:  market.RejectedNoOrders(mrktName),
	}
}

func newAccepted(
	price float64,
	priceInfo staticdb.PricingInfo,
) *ShopPrice {
	return &ShopPrice{
		Price: market.RoundedPrice(price),
		Desc: market.Accepted(
			priceInfo.MrktName,
			priceInfo.Prctile,
			priceInfo.Modifier,
			priceInfo.IsBuy,
		),
	}
}

func unpackPositivePrice(
	positivePrice *internal.PositivePrice,
	priceInfo staticdb.PricingInfo,
) *ShopPrice {
	accepted, price := market.UnpackPositivePrice(positivePrice)
	if accepted {
		return newAccepted(price, priceInfo)
	} else {
		return newRejectedNoOrders(priceInfo.MrktName)
	}
}
