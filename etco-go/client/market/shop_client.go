package market

import (
	"context"

	"github.com/WiggidyW/etco-go/client/market/marketprice"
)

type ShopPriceClient struct {
	marketPriceClient marketprice.MarketPriceClient
}

func NewShopPriceClient(
	marketPriceClient marketprice.MarketPriceClient,
) ShopPriceClient {
	return ShopPriceClient{marketPriceClient}
}

// gets the ShopPrice for the given item
// nil if price == 0 (instead of 0 with a rejection desc)
func (spc ShopPriceClient) Fetch(
	ctx context.Context,
	params ShopPriceParams,
) (*ShopPrice, error) {
	sPricingInfo := params.ShopLocationInfo.
		GetTypePricingInfo(params.TypeId)
	if sPricingInfo == nil { // item not sold at locations shop
		return newRejected(params.TypeId, params.Quantity), nil
	}

	pricePtr, err := spc.marketPriceClient.Fetch(
		ctx,
		marketprice.MarketPriceParams{
			PricingInfo: *sPricingInfo,
			TypeId:      params.TypeId,
		},
	)
	if err != nil {
		return nil, err
	}

	return unpackPositivePrice(
		params.TypeId,
		params.Quantity,
		pricePtr,
		*sPricingInfo,
	), nil
}
