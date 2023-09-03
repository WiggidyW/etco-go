package ordersregion

import (
	"context"

	pe "github.com/WiggidyW/etco-go/client/esi/model/internal/pageentries"
	pes "github.com/WiggidyW/etco-go/client/esi/model/internal/pageentries/streaming"
	"github.com/WiggidyW/etco-go/client/esi/raw_"
)

const ORDERS_REGION_ENTRIES_PER_PAGE int = 1000

type OrdersRegionClient struct {
	Inner pes.StreamingPageEntriesClient[
		OrdersRegionUrlParams,
		OrdersRegionEntry,
	]
}

func NewOrdersRegionClient(rawClient raw_.RawClient) OrdersRegionClient {
	return OrdersRegionClient{
		Inner: pes.NewStreamingPageEntriesClient[
			OrdersRegionUrlParams,
			OrdersRegionEntry,
		](
			rawClient,
			ORDERS_REGION_ENTRIES_PER_PAGE,
		),
	}
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
