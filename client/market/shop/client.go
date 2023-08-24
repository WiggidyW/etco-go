package shop

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/client/market/internal"
)

type ShopPriceClient struct {
	inner internal.MarketPriceClient
}

// gets the ShopPrice for the given item
// nil if price == 0 (instead of 0 with a rejection desc)
func (spc ShopPriceClient) Fetch(
	ctx context.Context,
	params ShopPriceParams,
) (*ShopPrice, error) {
	sTypeInfo := params.ShopLocationInfo.GetTypeInfo(params.TypeId)
	if sTypeInfo == nil { // item not sold at locations shop
		return newRejected(params.TypeId, params.Quantity), nil
	}

	pricePtr, err := spc.inner.Fetch(
		ctx,
		internal.MarketPriceParams{
			PricingInfo: *sTypeInfo,
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
		*sTypeInfo,
	), nil
}
