package raworders_

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	wc "github.com/WiggidyW/etco-go/client/caching/weak"
	mordersregion "github.com/WiggidyW/etco-go/client/esi/model/ordersregion"
)

const (
	RAW_REGION_MARKET_MIN_EXPIRES   time.Duration = 0
	RAW_REGION_MARKET_SLOCK_TTL     time.Duration = 20 * time.Second
	RAW_REGION_MARKET_SLOCK_MAXWAIT time.Duration = 10 * time.Second
)

type WC_RawRegionMarketClient = wc.WeakCachingClient[
	RawRegionMarketParams,
	RegionMarket,
	cache.ExpirableData[RegionMarket],
	RawRegionMarketClient,
]

func NewWC_RawRegionMarketClient(
	modelClient mordersregion.OrdersRegionClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) WC_RawRegionMarketClient {
	return wc.NewWeakCachingClient(
		NewRawRegionMarketClient(modelClient),
		RAW_REGION_MARKET_MIN_EXPIRES,
		cCache,
		sCache,
		RAW_REGION_MARKET_SLOCK_TTL,
		RAW_REGION_MARKET_SLOCK_MAXWAIT,
	)
}

type RawRegionMarketClient struct {
	modelClient mordersregion.OrdersRegionClient
}

func NewRawRegionMarketClient(
	modelClient mordersregion.OrdersRegionClient,
) RawRegionMarketClient {
	return RawRegionMarketClient{modelClient}
}

func (rmc RawRegionMarketClient) Fetch(
	ctx context.Context,
	params RawRegionMarketParams,
) (*cache.ExpirableData[RegionMarket], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the model receiver
	hrwc, err := rmc.modelClient.Fetch(
		ctx,
		mordersregion.OrdersRegionParams(params),
	)
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
			appendOrder(
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
) *[]MarketOrder {
	locationOrders, ok := m[locationId]
	if !ok {
		locationOrders = &[]MarketOrder{}
		m[locationId] = locationOrders
	}
	return locationOrders
}
