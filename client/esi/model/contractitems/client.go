package contractitems

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	e "github.com/WiggidyW/etco-go/client/esi/model/internal/entries"
	"github.com/WiggidyW/etco-go/client/esi/model/internal/naive"
	"github.com/WiggidyW/etco-go/client/esi/raw_"
)

const CONTRACT_ITEMS_ENTRIES_PER_PAGE int = 5000

type ContractItemsClient struct {
	Inner e.EntriesClient[
		ContractItemsUrlParams,
		ContractItemsEntry,
	]
}

func NewContractItemsClient(rawClient raw_.RawClient) ContractItemsClient {
	return ContractItemsClient{
		Inner: e.NewEntriesClient[
			ContractItemsUrlParams,
			ContractItemsEntry,
		](
			rawClient,
			CONTRACT_ITEMS_ENTRIES_PER_PAGE,
		),
	}
}

func (cic ContractItemsClient) Fetch(
	ctx context.Context,
	params ContractItemsParams,
) (*cache.ExpirableData[[]ContractItemsEntry], error) {
	return cic.Inner.Fetch(
		ctx,
		naive.NaiveParams[ContractItemsUrlParams]{
			UrlParams: ContractItemsUrlParams{
				CorporationId: params.CorporationId,
				ContractId:    params.ContractId,
			},
			AuthParams: &naive.AuthParams{
				Token: params.WebRefreshToken,
			},
		},
	)
}
