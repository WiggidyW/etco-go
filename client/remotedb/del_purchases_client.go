package remotedb

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	smac "github.com/WiggidyW/etco-go/client/caching/strong/anticaching/multi"
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

type SMAC_DelPurchasesClient = smac.StrongMultiAntiCachingClient[
	DelPurchasesParams,
	struct{},
	DelPurchasesClient,
]

func NewSMAC_DelPurchasesClient(
	rdbClient *rdb.RemoteDBClient,
	readShopQueueAntiCache *cache.StrongAntiCache,
	unreservedShopAssetsAntiCache *cache.StrongAntiCache,
) SMAC_DelPurchasesClient {
	return smac.NewStrongMultiAntiCachingClient(
		NewDelPurchasesClient(rdbClient),
		readShopQueueAntiCache,
		unreservedShopAssetsAntiCache,
	)
}

type DelPurchasesClient struct {
	rdbClient *rdb.RemoteDBClient
}

func NewDelPurchasesClient(rdbClient *rdb.RemoteDBClient) DelPurchasesClient {
	return DelPurchasesClient{rdbClient}
}

func (dpc DelPurchasesClient) Fetch(
	ctx context.Context,
	params DelPurchasesParams,
) (*struct{}, error) {
	if len(params.AppraisalCodes) == 0 {
		return &struct{}{}, nil
	}

	err := dpc.rdbClient.DelShopPurchases(ctx, params.AppraisalCodes...)

	if err != nil {
		return nil, err
	} else {
		return &struct{}{}, nil
	}
}
