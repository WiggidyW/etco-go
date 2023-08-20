package filterregion

import (
	"context"

	"github.com/WiggidyW/weve-esi/cache"
	wc "github.com/WiggidyW/weve-esi/client/caching/weak"
	"github.com/WiggidyW/weve-esi/client/market/internal/filter"
	"github.com/WiggidyW/weve-esi/client/market/internal/raw/region"
)

type WC_FilterRegionMarketClient = wc.WeakCachingClient[
	FilterRegionMarketParams,
	filter.SortedMarketOrders,
	cache.ExpirableData[filter.SortedMarketOrders],
	FilterRegionMarketClient,
]

type FilterRegionMarketClient struct {
	Inner region.WC_RegionMarketClient
}

// return is sorted by price, lowest first, and deduplicated
func (frmc FilterRegionMarketClient) Fetch(
	ctx context.Context,
	params FilterRegionMarketParams,
) (*cache.ExpirableData[filter.SortedMarketOrders], error) {
	marketRep, err := frmc.Inner.Fetch(
		ctx,
		region.RegionMarketParams{
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
			filter.SortedMarketOrders{},
			marketRep.Expires(),
		), nil
	}

	return cache.NewExpirableDataPtr(
		filter.SortDedupOrders(rawOrders),
		marketRep.Expires(),
	), nil
}
