package ordersregion

import (
	"context"

	pe "github.com/WiggidyW/eve-trading-co-go/client/esi/model/internal/pageentries"
	pes "github.com/WiggidyW/eve-trading-co-go/client/esi/model/internal/pageentries/streaming"
)

const ORDERS_REGION_ENTRIES_PER_PAGE int = 1000

type OrdersRegionClient struct {
	Inner pes.StreamingPageEntriesClient[
		OrdersRegionUrlParams,
		OrdersRegionEntry,
	]
}

func (orc OrdersRegionClient) Fetch(
	ctx context.Context,
	params OrdersRegionParams,
) (*pes.HeadRepWithChan[OrdersRegionEntry], error) {
	return orc.Inner.Fetch(
		ctx,
		pe.NaivePageParams[OrdersRegionUrlParams]{
			UrlParams: OrdersRegionUrlParams{
				RegionId:  params.RegionId,
				TypeId:    &params.TypeId,
				OrderType: boolToOrderType(params.IsBuy),
			},
		},
	)
}
