package filterorders_

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	wc "github.com/WiggidyW/etco-go/client/caching/weak"
	mordersregion "github.com/WiggidyW/etco-go/client/esi/model/ordersregion"
	"github.com/WiggidyW/etco-go/client/market/marketprice/filterorders_/raworders_"
)

const (
	FILTER_REGION_MARKET_MIN_EXPIRES   time.Duration = 0
	FILTER_REGION_MARKET_SLOCK_TTL     time.Duration = 20 * time.Second
	FILTER_REGION_MARKET_SLOCK_MAXWAIT time.Duration = 10 * time.Second
)

type WC_FilterRegionMarketClient = wc.WeakCachingClient[
	FilterRegionMarketParams,
	SortedMarketOrders,
	cache.ExpirableData[SortedMarketOrders],
	FilterRegionMarketClient,
]

func NewWC_FilterRegionMarketClient(
	modelClient mordersregion.OrdersRegionClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) WC_FilterRegionMarketClient {
	return wc.NewWeakCachingClient(
		NewFilterRegionMarketClient(
			modelClient,
			cCache,
			sCache,
		),
		FILTER_REGION_MARKET_MIN_EXPIRES,
		cCache,
		sCache,
		FILTER_REGION_MARKET_SLOCK_TTL,
		FILTER_REGION_MARKET_SLOCK_MAXWAIT,
	)
}

type FilterRegionMarketClient struct {
	rawRegionClient raworders_.WC_RawRegionMarketClient
}

func NewFilterRegionMarketClient(
	modelClient mordersregion.OrdersRegionClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) FilterRegionMarketClient {
	return FilterRegionMarketClient{
		raworders_.NewWC_RawRegionMarketClient(
			modelClient,
			cCache,
			sCache,
		),
	}
}

// return is sorted by price, lowest first, and deduplicated
func (frmc FilterRegionMarketClient) Fetch(
	ctx context.Context,
	params FilterRegionMarketParams,
) (*cache.ExpirableData[SortedMarketOrders], error) {
	marketRep, err := frmc.rawRegionClient.Fetch(
		ctx,
		raworders_.RawRegionMarketParams{
			RegionId: params.RegionId,
			TypeId:   params.TypeId,
			IsBuy:    params.IsBuy,
		},
	)
	if err != nil {
		return nil, err
	}

	// extract the orders that match the location
	rawOrders := marketRep.Data()[params.LocationId]
	if len(rawOrders) == 0 {
		// empty slice is worth caching
		return cache.NewExpirableDataPtr(
			SortedMarketOrders{},
			marketRep.Expires(),
		), nil
	}

	return cache.NewExpirableDataPtr(
		SortDedupOrders(rawOrders),
		marketRep.Expires(),
	), nil
}
