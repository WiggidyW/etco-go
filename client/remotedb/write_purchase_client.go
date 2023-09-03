package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	smac "github.com/WiggidyW/etco-go/client/caching/strong/anticaching/multi"
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

type SMAC_WritePurchaseClient = smac.StrongMultiAntiCachingClient[
	WritePurchaseParams,
	time.Time,
	WritePurchaseClient,
]

func NewSMAC_WritePurchaseClient(
	rdbClient *rdb.RemoteDBClient,
	readUserDataAntiCache *cache.StrongAntiCache,
	readShopQueueAntiCache *cache.StrongAntiCache,
	unreservedShopAssetsAntiCache *cache.StrongAntiCache,
) SMAC_WritePurchaseClient {
	return smac.NewStrongMultiAntiCachingClient(
		NewWritePurchaseClient(rdbClient),
		readUserDataAntiCache,
		readShopQueueAntiCache,
		unreservedShopAssetsAntiCache,
	)
}

type WritePurchaseClient struct {
	rdbClient *rdb.RemoteDBClient
}

func NewWritePurchaseClient(
	rdbClient *rdb.RemoteDBClient,
) WritePurchaseClient {
	return WritePurchaseClient{rdbClient}
}

// returns the time the appraisal was saved
func (wpc WritePurchaseClient) Fetch(
	ctx context.Context,
	params WritePurchaseParams,
) (*time.Time, error) {
	err := wpc.rdbClient.SaveShopPurchase(ctx, params.Appraisal)

	if err != nil {
		return nil, err
	} else {
		now := time.Now()
		return &now, nil
	}
}
