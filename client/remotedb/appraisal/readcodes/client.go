package readcodes

import (
	"context"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	sc "github.com/WiggidyW/weve-esi/client/caching/strong/caching"
	a "github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
)

type SC_ReadShopAppraisalCodesClient = sc.StrongCachingClient[
	ReadCharacterAppraisalCodesParams,
	a.CharacterAppraisalCodes,
	cache.ExpirableData[a.CharacterAppraisalCodes],
	ReadCharacterAppraisalCodesClient,
]

type ReadCharacterAppraisalCodesClient struct {
	Inner   *rdb.RemoteDBClient
	Expires time.Duration
}

func (rcacc ReadCharacterAppraisalCodesClient) Fetch(
	ctx context.Context,
	params ReadCharacterAppraisalCodesParams,
) (*cache.ExpirableData[a.CharacterAppraisalCodes], error) {
	codes, err := GetCharacterAppraisalCodes(
		ctx,
		rcacc.Inner,
		params.CharacterId,
	)
	if err != nil {
		return nil, err
	}

	return cache.NewExpirableDataPtr(
		*codes,
		time.Now().Add(rcacc.Expires),
	), nil
}
