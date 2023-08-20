package structureinfo

import (
	"context"

	"github.com/WiggidyW/weve-esi/cache"
	wc "github.com/WiggidyW/weve-esi/client/caching/weak"
	m "github.com/WiggidyW/weve-esi/client/esi/model/internal/model"
	"github.com/WiggidyW/weve-esi/client/esi/model/internal/naive"
)

type WC_StructureInfoClient = wc.WeakCachingClient[
	StructureInfoParams,
	StructureInfoModel,
	cache.ExpirableData[StructureInfoModel],
	StructureInfoClient,
]

type StructureInfoClient struct {
	Inner m.ModelClient[
		StructureInfoUrlParams,
		StructureInfoModel,
	]
}

func (sic StructureInfoClient) Fetch(
	ctx context.Context,
	params StructureInfoParams,
) (*cache.ExpirableData[StructureInfoModel], error) {
	return sic.Inner.Fetch(
		ctx,
		m.ModelParams[StructureInfoUrlParams, StructureInfoModel]{
			NaiveParams: naive.NaiveParams[StructureInfoUrlParams]{
				UrlParams: StructureInfoUrlParams{
					StructureId: params.StructureId,
				},
				AuthParams: &naive.AuthParams{
					Token: params.WebRefreshToken,
				},
			},
			Model: &StructureInfoModel{},
		},
	)
}
