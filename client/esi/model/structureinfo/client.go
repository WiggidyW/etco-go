package structureinfo

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	m "github.com/WiggidyW/etco-go/client/esi/model/internal/model"
	"github.com/WiggidyW/etco-go/client/esi/model/internal/naive"
	"github.com/WiggidyW/etco-go/client/esi/raw_"
)

type StructureInfoClient struct {
	Inner m.ModelClient[
		StructureInfoUrlParams,
		StructureInfoModel,
	]
}

func NewStructureInfoClient(rawClient raw_.RawClient) StructureInfoClient {
	return StructureInfoClient{
		Inner: m.NewModelClient[
			StructureInfoUrlParams,
			StructureInfoModel,
		](
			rawClient,
		),
	}
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
