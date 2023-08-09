package market

import (
	"context"
	"fmt"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/smartclient/market/percentile"
	"github.com/WiggidyW/weve-esi/logger"
	"github.com/WiggidyW/weve-esi/staticdb/tc"
)

const MIN_SHOP_PRICE = 0.01

type ShopClientFetchParams struct {
	TypeId     int32
	LocationId int64
}

type ShopClient struct {
	client *client.CachingClient[
		percentile.MarketPercentileClientFetchParams,
		percentile.MarketPercentile,
		cache.ExpirableData[percentile.MarketPercentile],
		*percentile.MarketPercentileClient,
	]
}

// gets the ShopPrice for the given item
// nil if price == 0 (instead of 0 with a rejection desc)
func (sc *ShopClient) Fetch(
	ctx context.Context,
	params ShopClientFetchParams,
) (*ShopPrice, error) {
	// static data
	// get the shop info
	shopInfo := tc.KVReaderShopInfo.Get(1) // capacity of 1
	// get the location
	location, ok := shopInfo.GetLocation(params.LocationId)
	if !ok { // location has no shop
		return nil, nil
	}
	// get the pricing
	pricing, ok := location.GetType(params.TypeId)
	if !ok { // type not sold at location
		return nil, nil
	}

	// validate the modifier
	modifier := pricing.Modifier()
	if modifier == 0 {
		logger.Logger.Fatal(fmt.Sprintf(
			"shop pricing modifier for %d at %d is 0",
			params.TypeId,
			params.LocationId,
		))
	}

	// fetch the percentile price
	percentile, err := sc.client.Fetch(
		ctx,
		percentile.NewFetchParams(pricing, params.TypeId),
	)
	if err != nil {
		return nil, err
	} else if percentile.Data().Price <= 0 {
		return nil, nil
	}

	// return the shop price
	shopPrice := newShopPrice(
		minPriced( // minned(rounded(multed))
			roundedToCents(multedByModifier(
				percentile.Data().Price,
				modifier,
			)),
			MIN_SHOP_PRICE,
		),
		newDesc(
			pricing.MarketName(),
			pricing.Percentile(),
			modifier,
			pricing.IsBuy(),
		),
	)
	return &shopPrice, nil
}
