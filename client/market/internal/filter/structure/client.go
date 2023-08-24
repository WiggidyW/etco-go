package filterstructure

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/cache"
	wc "github.com/WiggidyW/eve-trading-co-go/client/caching/weak"
	"github.com/WiggidyW/eve-trading-co-go/client/market/internal/filter"
	"github.com/WiggidyW/eve-trading-co-go/client/market/internal/raw"
	"github.com/WiggidyW/eve-trading-co-go/client/market/internal/raw/structure"
)

type WC_FilterStructureMarketClient = wc.WeakCachingClient[
	FilterStructureMarketParams,
	filter.SortedMarketOrders,
	cache.ExpirableData[filter.SortedMarketOrders],
	FilterStructureMarketClient,
]

type FilterStructureMarketClient struct {
	Inner structure.WC_StructureMarketClient
}

// return is sorted by price, lowest first, and deduplicated
func (fsmc FilterStructureMarketClient) Fetch(
	ctx context.Context,
	params FilterStructureMarketParams,
) (*cache.ExpirableData[filter.SortedMarketOrders], error) {
	marketRep, err := fsmc.Inner.Fetch(
		ctx,
		structure.StructureMarketParams{
			WebRefreshToken: params.WebRefreshToken,
			StructureId:     params.StructureId,
		},
	)
	if err != nil {
		return nil, err
	}

	// extract the orders that match the type and buy/sell
	var rawOrders []raw.MarketOrder
	allOrders := marketRep.Data()[params.TypeId]
	if params.IsBuy {
		rawOrders = allOrders.Buy
	} else {
		rawOrders = allOrders.Sell
	}
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
