package structureinfo

import (
	"context"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	wc "github.com/WiggidyW/etco-go/client/caching/weak"
	mstructureinfo "github.com/WiggidyW/etco-go/client/esi/model/structureinfo"
)

const (
	STRUCTURE_INFO_MIN_EXPIRES   time.Duration = 48 * time.Hour
	STRUCTURE_INFO_SLOCK_TTL     time.Duration = 30 * time.Second
	STRUCTURE_INFO_SLOCK_MAXWAIT time.Duration = 10 * time.Second
)

type WC_StructureInfoClient = wc.WeakCachingClient[
	StructureInfoParams,
	StructureInfo,
	cache.ExpirableData[StructureInfo],
	StructureInfoClient,
]

func NewWC_StructureInfoClient(
	modelClient mstructureinfo.StructureInfoClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) WC_StructureInfoClient {
	return wc.NewWeakCachingClient(
		NewStructureInfoClient(modelClient),
		STRUCTURE_INFO_MIN_EXPIRES,
		cCache,
		sCache,
		STRUCTURE_INFO_SLOCK_TTL,
		STRUCTURE_INFO_SLOCK_MAXWAIT,
	)
}

type StructureInfoClient struct {
	modelClient mstructureinfo.StructureInfoClient
}

func NewStructureInfoClient(
	modelClient mstructureinfo.StructureInfoClient,
) StructureInfoClient {
	return StructureInfoClient{modelClient}
}

func (sic StructureInfoClient) Fetch(
	ctx context.Context,
	params StructureInfoParams,
) (*cache.ExpirableData[StructureInfo], error) {
	modelRep, err := sic.modelClient.Fetch(
		ctx,
		mstructureinfo.StructureInfoParams{
			WebRefreshToken: build.STRUCTURE_INFO_WEB_REFRESH_TOKEN,
			StructureId:     params.StructureId,
		},
	)

	if err != nil && Forbidden(err) {
		return cache.NewExpirableDataPtr(
			StructureInfo{
				Forbidden: true,
				// Name:      "",
				SystemId: -1,
			},
			time.Time{},
		), nil

	} else if err != nil {
		return nil, err

	} else {
		return cache.NewExpirableDataPtr(
			StructureInfo{
				Forbidden: false,
				Name:      modelRep.Data().Name,
				SystemId:  modelRep.Data().SolarSystemId,
			},
			modelRep.Expires(),
		), nil
	}

}
