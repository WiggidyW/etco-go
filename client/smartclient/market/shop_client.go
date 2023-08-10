package market

import (
	"context"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/smartclient/market/desc"
	"github.com/WiggidyW/weve-esi/client/smartclient/market/percentile"
	"github.com/WiggidyW/weve-esi/client/smartclient/market/percentile/orders"
	"github.com/WiggidyW/weve-esi/staticdb"
)

const MIN_SHOP_PRICE = 0.01 // minimum price for non-0-base-price items

type ShopParams struct {
	TypeId     int32
	LocationId int64
}

type ShopClient struct {
	client *client.CachingClient[
		percentile.MrktPrctileParams,
		percentile.MrktPrctile,
		cache.ExpirableData[percentile.MrktPrctile],
		*percentile.MrktPrctileClient,
	]
}

// gets the ShopPrice for the given item
// nil if price == 0 (instead of 0 with a rejection desc)
func (sc ShopClient) Fetch(
	ctx context.Context,
	params ShopParams,
) (*ShopPrice, error) {
	// // static data
	// get the location
	shopLocation := staticdb.GetShopLocationInfo(params.LocationId)
	if shopLocation == nil { // location has no shop
		return nil, nil
	}
	// get the type info / pricing
	shopTypeInfo := shopLocation.GetTypeInfo(params.TypeId)
	if shopTypeInfo == nil { // item not sold at locations shop
		return nil, nil
	}

	// fetch the percentile price
	var prctile percentile.MrktPrctile
	if prctileRep, err := sc.client.Fetch(
		ctx,
		percentile.MrktPrctileParams{
			MrktOrdersParams: orders.MrktOrdersParams{
				PricingInfo: *shopTypeInfo,
				TypeId:      params.TypeId,
			},
		},
	); err != nil {
		return nil, err
	} else {
		prctile = prctileRep.Data()
	}

	// if the price is 0, return nil
	if prctile.Price <= 0 {
		return nil, nil
	}

	// return the shop price, it will be at least MIN_SHOP_PRICE
	shopPrice := newShopPrice(
		minPriced( // minned(rounded(multed))
			roundedToCents(multedByModifier(
				prctile.Price,
				shopTypeInfo.Modifier,
			)),
			MIN_SHOP_PRICE,
		),
		desc.Accepted(
			shopTypeInfo.MrktName,
			shopTypeInfo.Prctile,
			shopTypeInfo.Modifier,
			shopTypeInfo.IsBuy,
		),
	)
	return &shopPrice, nil
}
