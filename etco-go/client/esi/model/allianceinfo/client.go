package allianceinfo

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
	ALLIANCE_INFO_MIN_EXPIRES    time.Duration = 24 * time.Hour
	ALLIANCE_INFO_SLOCK_TTL      time.Duration = 30 * time.Second
	ALLIANCE_INFO_SLOCK_MAX_WAIT time.Duration = 10 * time.Second
)

type WC_AllianceInfoClient = wc.WeakCachingClient[
	AllianceInfoParams,
	AllianceInfoModel,
	cache.ExpirableData[AllianceInfoModel],
	AllianceInfoClient,
]

func NewWC_AllianceInfoClient(
	rawClient raw_.RawClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) WC_AllianceInfoClient {
	return wc.NewWeakCachingClient[
		AllianceInfoParams,
		AllianceInfoModel,
		cache.ExpirableData[AllianceInfoModel],
		AllianceInfoClient,
	](
		NewAllianceInfoClient(rawClient),
		ALLIANCE_INFO_MIN_EXPIRES,
		cCache,
		sCache,
		ALLIANCE_INFO_SLOCK_TTL,
		ALLIANCE_INFO_SLOCK_MAX_WAIT,
	)
}

type AllianceInfoClient struct {
	Inner m.ModelClient[
		AllianceInfoUrlParams,
		AllianceInfoModel,
	]
}

func NewAllianceInfoClient(rawClient raw_.RawClient) AllianceInfoClient {
	return AllianceInfoClient{
		Inner: m.NewModelClient[
			AllianceInfoUrlParams,
			AllianceInfoModel,
		](
			rawClient,
		),
	}
}

func (aic AllianceInfoClient) Fetch(
	ctx context.Context,
	params AllianceInfoParams,
) (*cache.ExpirableData[AllianceInfoModel], error) {
	return aic.Inner.Fetch(
		ctx,
		m.ModelParams[AllianceInfoUrlParams, AllianceInfoModel]{
			NaiveParams: naive.NaiveParams[AllianceInfoUrlParams]{
				UrlParams: AllianceInfoUrlParams(params),
			},
			Model: &AllianceInfoModel{},
		},
	)
}
