package characterinfo

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	wc "github.com/WiggidyW/etco-go/client/caching/weak"
	m "github.com/WiggidyW/etco-go/client/esi/model/internal/model"
	"github.com/WiggidyW/etco-go/client/esi/model/internal/naive"
	"github.com/WiggidyW/etco-go/client/esi/raw_"
)

const (
	CHARACTER_INFO_MIN_EXPIRES    time.Duration = 24 * time.Hour
	CHARACTER_INFO_SLOCK_TTL      time.Duration = 30 * time.Second
	CHARACTER_INFO_SLOCK_MAX_WAIT time.Duration = 10 * time.Second
)

type WC_CharacterInfoClient = wc.WeakCachingClient[
	CharacterInfoParams,
	CharacterInfoModel,
	cache.ExpirableData[CharacterInfoModel],
	CharacterInfoClient,
]

func NewWC_CharacterInfoClient(
	rawClient raw_.RawClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) WC_CharacterInfoClient {
	return wc.NewWeakCachingClient[
		CharacterInfoParams,
		CharacterInfoModel,
		cache.ExpirableData[CharacterInfoModel],
		CharacterInfoClient,
	](
		NewCharacterInfoClient(rawClient),
		CHARACTER_INFO_MIN_EXPIRES,
		cCache,
		sCache,
		CHARACTER_INFO_SLOCK_TTL,
		CHARACTER_INFO_SLOCK_MAX_WAIT,
	)
}

type CharacterInfoClient struct {
	Inner m.ModelClient[
		CharacterInfoUrlParams,
		CharacterInfoModel,
	]
}

func NewCharacterInfoClient(rawClient raw_.RawClient) CharacterInfoClient {
	return CharacterInfoClient{
		Inner: m.NewModelClient[
			CharacterInfoUrlParams,
			CharacterInfoModel,
		](
			rawClient,
		),
	}
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
