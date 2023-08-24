package contractscorporation

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/client/esi/model/internal/naive"
	pe "github.com/WiggidyW/eve-trading-co-go/client/esi/model/internal/pageentries"
	pes "github.com/WiggidyW/eve-trading-co-go/client/esi/model/internal/pageentries/streaming"
)

const CONTRACTS_CORPORATION_ENTRIES_PER_PAGE int = 1000

type ContractsCorporationClient struct {
	Inner pes.StreamingPageEntriesClient[
		ContractsCorporationUrlParams,
		ContractsCorporationEntry,
	]
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
