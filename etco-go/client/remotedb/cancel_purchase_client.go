package remotedb

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	smac "github.com/WiggidyW/etco-go/client/caching/strong/anticaching/multi"
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

type SMAC_CancelPurchaseClient = smac.StrongMultiAntiCachingClient[
	CancelPurchaseParams,
	struct{},
	CancelPurchaseClient,
]

func NewSMAC_CancelPurchaseClient(
	rdbClient *rdb.RemoteDBClient,
	readUserDataAntiCache *cache.StrongAntiCache,
	readShopQueueAntiCache *cache.StrongAntiCache,
	unreservedShopAssetsAntiCache *cache.StrongAntiCache,
) SMAC_CancelPurchaseClient {
	return smac.NewStrongMultiAntiCachingClient(
		NewCancelPurchaseClient(rdbClient),
		readUserDataAntiCache,
		readShopQueueAntiCache,
		unreservedShopAssetsAntiCache,
	)
}

type CancelPurchaseClient struct {
	rdbClient *rdb.RemoteDBClient
}

func NewCancelPurchaseClient(rdbClient *rdb.RemoteDBClient) CancelPurchaseClient {
	return CancelPurchaseClient{rdbClient}
}

// returns the time the appraisal was saved
func (cpc CancelPurchaseClient) Fetch(
	ctx context.Context,
	params CancelPurchaseParams,
) (*struct{}, error) {
	err := cpc.rdbClient.CancelShopPurchase(
		ctx,
		params.CharacterId,
		params.AppraisalCode,
	)

	if err != nil {
		return nil, err
	} else {
		return &struct{}{}, nil
	}
}
