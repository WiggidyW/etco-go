package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	sc "github.com/WiggidyW/etco-go/client/caching/strong/caching"
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

const (
	READ_USER_DATA_EXPIRES        time.Duration = 24 * time.Hour
	READ_USER_DATA_MIN_EXPIRES    time.Duration = 0
	READ_USER_DATA_SLOCK_TTL      time.Duration = 1 * time.Minute
	READ_USER_DATA_SLOCK_MAX_WAIT time.Duration = 1 * time.Minute
)

type SC_ReadUserDataClient = sc.StrongCachingClient[
	ReadUserDataParams,
	rdb.UserData,
	cache.ExpirableData[rdb.UserData],
	ReadUserDataClient,
]

func NewSC_ReadUserDataClient(
	rdbClient *rdb.RemoteDBClient,
	sCache cache.SharedServerCache,
) SC_ReadUserDataClient {
	return sc.NewStrongCachingClient(
		NewReadUserDataClient(rdbClient),
		READ_USER_DATA_MIN_EXPIRES,
		sCache,
		READ_USER_DATA_SLOCK_TTL,
		READ_USER_DATA_SLOCK_MAX_WAIT,
	)
}

type ReadUserDataClient struct {
	rdbClient *rdb.RemoteDBClient
	expires   time.Duration
}

func NewReadUserDataClient(
	rdbClient *rdb.RemoteDBClient,
) ReadUserDataClient {
	return ReadUserDataClient{
		rdbClient: rdbClient,
		expires:   READ_USER_DATA_EXPIRES,
	}
}

func (rudc ReadUserDataClient) Fetch(
	ctx context.Context,
	params ReadUserDataParams,
) (*cache.ExpirableData[rdb.UserData], error) {
	ud, err := rudc.rdbClient.ReadUserData(ctx, params.CharacterId)

	if err != nil {
		return nil, err
	} else {
		return cache.NewExpirableDataPtr(
			ud,
			time.Now().Add(rudc.expires),
		), nil
	}
}
