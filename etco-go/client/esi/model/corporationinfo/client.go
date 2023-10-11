package corporationinfo

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
	CORPORATION_INFO_MIN_EXPIRES    time.Duration = 24 * time.Hour
	CORPORATION_INFO_SLOCK_TTL      time.Duration = 30 * time.Second
	CORPORATION_INFO_SLOCK_MAX_WAIT time.Duration = 10 * time.Second
)

type WC_CorporationInfoClient = wc.WeakCachingClient[
	CorporationInfoParams,
	CorporationInfoModel,
	cache.ExpirableData[CorporationInfoModel],
	CorporationInfoClient,
]

func NewWC_CorporationInfoClient(
	rawClient raw_.RawClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) WC_CorporationInfoClient {
	return wc.NewWeakCachingClient[
		CorporationInfoParams,
		CorporationInfoModel,
		cache.ExpirableData[CorporationInfoModel],
		CorporationInfoClient,
	](
		NewCorporationInfoClient(rawClient),
		CORPORATION_INFO_MIN_EXPIRES,
		cCache,
		sCache,
		CORPORATION_INFO_SLOCK_TTL,
		CORPORATION_INFO_SLOCK_MAX_WAIT,
	)
}

type CorporationInfoClient struct {
	Inner m.ModelClient[
		CorporationInfoUrlParams,
		CorporationInfoModel,
	]
}

func NewCorporationInfoClient(rawClient raw_.RawClient) CorporationInfoClient {
	return CorporationInfoClient{
		Inner: m.NewModelClient[
			CorporationInfoUrlParams,
			CorporationInfoModel,
		](
			rawClient,
		),
	}
}

func (cic CorporationInfoClient) Fetch(
	ctx context.Context,
	params CorporationInfoParams,
) (*cache.ExpirableData[CorporationInfoModel], error) {
	return cic.Inner.Fetch(
		ctx,
		m.ModelParams[CorporationInfoUrlParams, CorporationInfoModel]{
			NaiveParams: naive.NaiveParams[CorporationInfoUrlParams]{
				UrlParams: CorporationInfoUrlParams(params),
			},
			Model: &CorporationInfoModel{},
		},
	)
}
