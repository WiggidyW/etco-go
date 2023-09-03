package assetscorporation

import (
	"context"

	"github.com/WiggidyW/etco-go/client/esi/model/internal/naive"
	pe "github.com/WiggidyW/etco-go/client/esi/model/internal/pageentries"
	pes "github.com/WiggidyW/etco-go/client/esi/model/internal/pageentries/streaming"
	"github.com/WiggidyW/etco-go/client/esi/raw_"
)

const ASSETS_CORPORATION_ENTRIES_PER_PAGE int = 1000

type AssetsCorporationClient struct {
	Inner pes.StreamingPageEntriesClient[
		AssetsCorporationUrlParams,
		AssetsCorporationEntry,
	]
}

func NewAssetsCorporationClient(
	rawClient raw_.RawClient,
) AssetsCorporationClient {
	return AssetsCorporationClient{
		Inner: pes.NewStreamingPageEntriesClient[
			AssetsCorporationUrlParams,
			AssetsCorporationEntry,
		](
			rawClient,
			ASSETS_CORPORATION_ENTRIES_PER_PAGE,
		),
	}
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
