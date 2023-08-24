package removematching

import (
	"context"

	smac "github.com/WiggidyW/eve-trading-co-go/client/caching/strong/anticaching/multi"
	rdb "github.com/WiggidyW/eve-trading-co-go/client/remotedb/internal"
)

type SMAC_ShopQueueRemoveMatchingClient = smac.StrongMultiAntiCachingClient[
	ShopQueueRemoveMatchingParams,
	struct{},
	ShopQueueRemoveMatchingClient,
]

type ShopQueueRemoveMatchingClient struct {
	Inner *rdb.RemoteDBClient
}

func (sqrmc ShopQueueRemoveMatchingClient) Fetch(
	ctx context.Context,
	params ShopQueueRemoveMatchingParams,
) (*struct{}, error) {
	err := SetShopQueueRemoveMatching(sqrmc.Inner, ctx, []string(params))
	if err != nil {
		return nil, err
	}
	return &struct{}{}, nil
}
