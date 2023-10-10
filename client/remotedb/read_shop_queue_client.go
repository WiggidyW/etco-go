package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	sc "github.com/WiggidyW/etco-go/client/caching/strong/caching"
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

const (
	READ_SHOP_QUEUE_EXPIRES        time.Duration = 24 * time.Hour
	READ_SHOP_QUEUE_MIN_EXPIRES    time.Duration = 0
	READ_SHOP_QUEUE_SLOCK_TTL      time.Duration = 1 * time.Minute
	READ_SHOP_QUEUE_SLOCK_MAX_WAIT time.Duration = 1 * time.Minute
)

type SC_ReadShopQueueClient = sc.StrongCachingClient[
	ReadShopQueueParams,
	[]string,
	cache.ExpirableData[[]string],
	ReadShopQueueClient,
]

func NewSC_ReadShopQueueClient(
	rdbClient *rdb.RemoteDBClient,
	sCache cache.SharedServerCache,
) SC_ReadShopQueueClient {
	return sc.NewStrongCachingClient(
		NewReadShopQueueClient(rdbClient),
		READ_SHOP_QUEUE_MIN_EXPIRES,
		sCache,
		READ_SHOP_QUEUE_SLOCK_TTL,
		READ_SHOP_QUEUE_SLOCK_MAX_WAIT,
	)
}

type ReadShopQueueClient struct {
	rdbClient *rdb.RemoteDBClient
	expires   time.Duration
}

func NewReadShopQueueClient(
	rdbClient *rdb.RemoteDBClient,
) ReadShopQueueClient {
	return ReadShopQueueClient{
		rdbClient: rdbClient,
		expires:   READ_SHOP_QUEUE_EXPIRES,
	}
}

func (sqrc ReadShopQueueClient) Fetch(
	ctx context.Context,
	params ReadShopQueueParams,
) (*cache.ExpirableData[[]string], error) {
	sq, err := sqrc.rdbClient.ReadShopQueue(ctx)

	if err != nil {
		return nil, err
	} else {
		return cache.NewExpirableDataPtr[[]string](
			sq.ShopQueue,
			time.Now().Add(sqrc.expires),
		), nil
	}
}
