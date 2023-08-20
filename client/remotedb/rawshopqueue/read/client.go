package read

import (
	"context"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	sc "github.com/WiggidyW/weve-esi/client/caching/strong/caching"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
)

type SC_ShopQueueReadClient = sc.StrongCachingClient[
	ShopQueueReadParams,
	[]string,
	cache.ExpirableData[[]string],
	ShopQueueReadClient,
]

type ShopQueueReadClient struct {
	Inner   *rdb.RemoteDBClient
	Expires time.Duration
}

func (sqrc ShopQueueReadClient) Fetch(
	ctx context.Context,
	params ShopQueueReadParams,
) (*cache.ExpirableData[[]string], error) {
	sq, err := GetShopQueue(sqrc.Inner, ctx)
	if err != nil {
		return nil, err
	}

	return cache.NewExpirableDataPtr[[]string](
		sq,
		time.Now().Add(sqrc.Expires),
	), nil
}
