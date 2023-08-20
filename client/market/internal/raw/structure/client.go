package structure

import (
	"context"

	"github.com/WiggidyW/weve-esi/cache"
	wc "github.com/WiggidyW/weve-esi/client/caching/weak"
	os "github.com/WiggidyW/weve-esi/client/esi/model/ordersstructure"
	"github.com/WiggidyW/weve-esi/client/market/internal/raw"
)

type WC_StructureMarketClient = wc.WeakCachingClient[
	StructureMarketParams,
	StructureMarket,
	cache.ExpirableData[StructureMarket],
	StructureMarketClient,
]

type StructureMarketClient struct {
	Inner os.OrdersStructureClient
}

func (smc StructureMarketClient) Fetch(
	ctx context.Context,
	params StructureMarketParams,
) (*cache.ExpirableData[StructureMarket], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the model receiver
	hrwc, err := smc.Inner.Fetch(ctx, os.OrdersStructureParams(params))
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
			raw.AppendOrder(
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
) *[]raw.MarketOrder {
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
