package region

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/cache"
	wc "github.com/WiggidyW/eve-trading-co-go/client/caching/weak"
	or "github.com/WiggidyW/eve-trading-co-go/client/esi/model/ordersregion"
	"github.com/WiggidyW/eve-trading-co-go/client/market/internal/raw"
)

type WC_RegionMarketClient = wc.WeakCachingClient[
	RegionMarketParams,
	RegionMarket,
	cache.ExpirableData[RegionMarket],
	RegionMarketClient,
]

type RegionMarketClient struct {
	Inner or.OrdersRegionClient
}

func (rmc RegionMarketClient) Fetch(
	ctx context.Context,
	params RegionMarketParams,
) (*cache.ExpirableData[RegionMarket], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the model receiver
	hrwc, err := rmc.Inner.Fetch(ctx, or.OrdersRegionParams(params))
	if err != nil {
		return nil, err
	}

	// // insert all entries into a map, keyed by location id
	regionMarket := make(initRegionMarket)
	for i := 0; i < hrwc.NumPages; i++ {

		// receive the next page
		page, err := hrwc.RecvUpdateExpires()
		if err != nil {
			return nil, err
		}

		for _, entry := range page {
			// append the order to the location orders
			locationOrders := getLocationOrders(
				regionMarket,
				entry.LocationId,
			)
			raw.AppendOrder(
				locationOrders,
				entry.Price,
				int64(entry.VolumeRemain),
			)
		}
	} // //

	return cache.NewExpirableDataPtr(
		finishRegionMarket(regionMarket), // *[] -> []
		hrwc.Expires,
	), nil
}

func getLocationOrders(
	m initRegionMarket,
	locationId int64,
) *[]raw.MarketOrder {
	locationOrders, ok := m[locationId]
	if !ok {
		locationOrders = &[]raw.MarketOrder{}
		m[locationId] = locationOrders
	}
	return locationOrders
}
