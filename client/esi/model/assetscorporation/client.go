package assetscorporation

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/esi/model/internal/naive"
	pe "github.com/WiggidyW/weve-esi/client/esi/model/internal/pageentries"
	pes "github.com/WiggidyW/weve-esi/client/esi/model/internal/pageentries/streaming"
)

const ASSETS_CORPORATION_ENTRIES_PER_PAGE int = 1000

type AssetsCorporationClient struct {
	Inner pes.StreamingPageEntriesClient[
		AssetsCorporationUrlParams,
		AssetsCorporationEntry,
	]
}

func (acc AssetsCorporationClient) Fetch(
	ctx context.Context,
	params AssetsCorporationParams,
) (*pes.HeadRepWithChan[AssetsCorporationEntry], error) {
	return acc.Inner.Fetch(
		ctx,
		pe.NaivePageParams[AssetsCorporationUrlParams]{
			UrlParams: AssetsCorporationUrlParams{
				CorporationId: params.CorporationId,
			},
			AuthParams: &naive.AuthParams{
				Token: params.WebRefreshToken,
			},
		},
	)
}
