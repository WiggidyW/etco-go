package characterinfo

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/cache"
	wc "github.com/WiggidyW/eve-trading-co-go/client/caching/weak"
	m "github.com/WiggidyW/eve-trading-co-go/client/esi/model/internal/model"
	"github.com/WiggidyW/eve-trading-co-go/client/esi/model/internal/naive"
)

type WC_CharacterInfoClient = wc.WeakCachingClient[
	CharacterInfoParams,
	CharacterInfoModel,
	cache.ExpirableData[CharacterInfoModel],
	CharacterInfoClient,
]

type CharacterInfoClient struct {
	Inner m.ModelClient[
		CharacterInfoUrlParams,
		CharacterInfoModel,
	]
}

func (cic CharacterInfoClient) Fetch(
	ctx context.Context,
	params CharacterInfoParams,
) (*cache.ExpirableData[CharacterInfoModel], error) {
	return cic.Inner.Fetch(
		ctx,
		m.ModelParams[CharacterInfoUrlParams, CharacterInfoModel]{
			NaiveParams: naive.NaiveParams[CharacterInfoUrlParams]{
				UrlParams: CharacterInfoUrlParams(params),
			},
			Model: &CharacterInfoModel{},
		},
	)
}
