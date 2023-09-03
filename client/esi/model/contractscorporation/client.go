package contractscorporation

import (
	"context"

	"github.com/WiggidyW/etco-go/client/esi/model/internal/naive"
	pe "github.com/WiggidyW/etco-go/client/esi/model/internal/pageentries"
	pes "github.com/WiggidyW/etco-go/client/esi/model/internal/pageentries/streaming"
	"github.com/WiggidyW/etco-go/client/esi/raw_"
)

const CONTRACTS_CORPORATION_ENTRIES_PER_PAGE int = 1000

type ContractsCorporationClient struct {
	Inner pes.StreamingPageEntriesClient[
		ContractsCorporationUrlParams,
		ContractsCorporationEntry,
	]
}

func NewContractsCorporationClient(
	rawClient raw_.RawClient,
) ContractsCorporationClient {
	return ContractsCorporationClient{
		Inner: pes.NewStreamingPageEntriesClient[
			ContractsCorporationUrlParams,
			ContractsCorporationEntry,
		](
			rawClient,
			CONTRACTS_CORPORATION_ENTRIES_PER_PAGE,
		),
	}
}

func (ccc ContractsCorporationClient) Fetch(
	ctx context.Context,
	params ContractsCorporationParams,
) (*pes.HeadRepWithChan[ContractsCorporationEntry], error) {
	return ccc.Inner.Fetch(
		ctx,
		pe.NaivePageParams[ContractsCorporationUrlParams]{
			UrlParams: ContractsCorporationUrlParams{
				CorporationId: params.CorporationId,
			},
			AuthParams: &naive.AuthParams{
				Token: params.WebRefreshToken,
			},
		},
	)
}
