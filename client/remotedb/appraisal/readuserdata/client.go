package readuserdata

import (
	"context"
	"time"

	"github.com/WiggidyW/eve-trading-co-go/cache"
	sc "github.com/WiggidyW/eve-trading-co-go/client/caching/strong/caching"
	a "github.com/WiggidyW/eve-trading-co-go/client/remotedb/appraisal"
	rdb "github.com/WiggidyW/eve-trading-co-go/client/remotedb/internal"
)

type SC_ReadUserDataClient = sc.StrongCachingClient[
	ReadUserDataParams,
	a.UserData,
	cache.ExpirableData[a.UserData],
	ReadUserDataClient,
]

type ReadUserDataClient struct {
	Inner   *rdb.RemoteDBClient
	Expires time.Duration
}

func (rcacc ReadUserDataClient) Fetch(
	ctx context.Context,
	params ReadUserDataParams,
) (*cache.ExpirableData[a.UserData], error) {
	uData, err := GetUserData(
		ctx,
		rcacc.Inner,
		params.CharacterId,
	)
	if err != nil {
		return nil, err
	}

	return cache.NewExpirableDataPtr(
		*uData,
		time.Now().Add(rcacc.Expires),
	), nil
}
