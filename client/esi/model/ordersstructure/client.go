package ordersstructure

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/esi/model/internal/naive"
	pe "github.com/WiggidyW/weve-esi/client/esi/model/internal/pageentries"
	pes "github.com/WiggidyW/weve-esi/client/esi/model/internal/pageentries/streaming"
)

const ORDERS_STRUCTURE_ENTRIES_PER_PAGE int = 1000

type OrdersStructureClient struct {
	Inner pes.StreamingPageEntriesClient[
		OrdersStructureUrlParams,
		OrdersStructureEntry,
	]
}

func (osc OrdersStructureClient) Fetch(
	ctx context.Context,
	params OrdersStructureParams,
) (*pes.HeadRepWithChan[OrdersStructureEntry], error) {
	return osc.Inner.Fetch(
		ctx,
		pe.NaivePageParams[OrdersStructureUrlParams]{
			UrlParams: OrdersStructureUrlParams{
				StructureId: params.StructureId,
			},
			AuthParams: &naive.AuthParams{
				Token: params.WebRefreshToken,
			},
		},
	)
}
