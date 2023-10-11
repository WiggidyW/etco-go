package filterorders_

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	wc "github.com/WiggidyW/etco-go/client/caching/weak"
	mordersstructure "github.com/WiggidyW/etco-go/client/esi/model/ordersstructure"
	"github.com/WiggidyW/etco-go/client/market/marketprice/filterorders_/raworders_"
)

const (
	FILTER_STRUCTURE_MARKET_MIN_EXPIRES   time.Duration = 0
	FILTER_STRUCTURE_MARKET_SLOCK_TTL     time.Duration = 30 * time.Second
	FILTER_STRUCTURE_MARKET_SLOCK_MAXWAIT time.Duration = 10 * time.Second
)

type WC_FilterStructureMarketClient = wc.WeakCachingClient[
	FilterStructureMarketParams,
	SortedMarketOrders,
	cache.ExpirableData[SortedMarketOrders],
	FilterStructureMarketClient,
]

func NewWC_FilterStructureMarketClient(
	modelClient mordersstructure.OrdersStructureClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) WC_FilterStructureMarketClient {
	return wc.NewWeakCachingClient(
		NewFilterStructureMarketClient(
			modelClient,
			cCache,
			sCache,
		),
		FILTER_STRUCTURE_MARKET_MIN_EXPIRES,
		cCache,
		sCache,
		FILTER_STRUCTURE_MARKET_SLOCK_TTL,
		FILTER_STRUCTURE_MARKET_SLOCK_MAXWAIT,
	)
}

type FilterStructureMarketClient struct {
	rawStructureClient raworders_.WC_RawStructureMarketClient
}

func NewFilterStructureMarketClient(
	modelClient mordersstructure.OrdersStructureClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) FilterStructureMarketClient {
	return FilterStructureMarketClient{
		raworders_.NewWC_RawStructureMarketClient(
			modelClient,
			cCache,
			sCache,
		),
	}
}

// return is sorted by price, lowest first, and deduplicated
func (fsmc FilterStructureMarketClient) Fetch(
	ctx context.Context,
	params FilterStructureMarketParams,
) (*cache.ExpirableData[SortedMarketOrders], error) {
	marketRep, err := fsmc.rawStructureClient.Fetch(
		ctx,
		raworders_.RawStructureMarketParams{
			WebRefreshToken: params.WebRefreshToken,
			StructureId:     params.StructureId,
		},
	)
	if err != nil {
		return nil, err
	}

	// extract the orders that match the type and buy/sell
	var rawOrders []raworders_.MarketOrder
	allOrders := marketRep.Data()[params.TypeId]
	if params.IsBuy {
		rawOrders = allOrders.Buy
	} else {
		rawOrders = allOrders.Sell
	}
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
