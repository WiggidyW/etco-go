package raworders_

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	wc "github.com/WiggidyW/etco-go/client/caching/weak"
	mordersstructure "github.com/WiggidyW/etco-go/client/esi/model/ordersstructure"
)

const (
	RAW_STRUCTURE_MARKET_MIN_EXPIRES   time.Duration = 30 * time.Minute
	RAW_STRUCTURE_MARKET_SLOCK_TTL     time.Duration = 1 * time.Minute
	RAW_STRUCTURE_MARKET_SLOCK_MAXWAIT time.Duration = 30 * time.Second
)

type WC_RawStructureMarketClient = wc.WeakCachingClient[
	RawStructureMarketParams,
	StructureMarket,
	cache.ExpirableData[StructureMarket],
	RawStructureMarketClient,
]

func NewWC_RawStructureMarketClient(
	modelClient mordersstructure.OrdersStructureClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) WC_RawStructureMarketClient {
	return wc.NewWeakCachingClient(
		NewRawStructureMarketClient(modelClient),
		RAW_STRUCTURE_MARKET_MIN_EXPIRES,
		cCache,
		sCache,
		RAW_STRUCTURE_MARKET_SLOCK_TTL,
		RAW_STRUCTURE_MARKET_SLOCK_MAXWAIT,
	)
}

type RawStructureMarketClient struct {
	modelClient mordersstructure.OrdersStructureClient
}

func NewRawStructureMarketClient(
	modelClient mordersstructure.OrdersStructureClient,
) RawStructureMarketClient {
	return RawStructureMarketClient{modelClient}
}

func (smc RawStructureMarketClient) Fetch(
	ctx context.Context,
	params RawStructureMarketParams,
) (*cache.ExpirableData[StructureMarket], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the model receiver
	hrwc, err := smc.modelClient.Fetch(
		ctx,
		mordersstructure.OrdersStructureParams(params),
	)
	if err != nil {
		return nil, err
	}

	// // insert all entries into a map, keyed by type id
	structureMarket := make(initStructureMarket)
	for i := 0; i < hrwc.NumPages; i++ {

		// receive the next page
		page, err := hrwc.RecvUpdateExpires()
		if err != nil {
			return nil, err
		}

		for _, entry := range page {
			// append the order to the type orders
			typeOrders := getTypeOrders(
				structureMarket,
				entry.TypeId,
				entry.IsBuyOrder,
			)
			appendOrder(
				typeOrders,
				entry.Price,
				int64(entry.VolumeRemain),
			)
		}
	} // //

	return cache.NewExpirableDataPtr(
		finishStructureMarket(structureMarket), // *T -> T
		hrwc.Expires,
	), nil
}

func getTypeOrders(
	m initStructureMarket,
	typeId int32,
	isBuy bool,
) *[]MarketOrder {
	typeOrders, ok := m[typeId]
	if !ok {
		typeOrders = &StructureOrders{}
		m[typeId] = typeOrders
	}

	if isBuy {
		return &typeOrders.Buy
	} else {
		return &typeOrders.Sell
	}
}
