package write

import (
	"context"
	"time"

	smac "github.com/WiggidyW/eve-trading-co-go/client/caching/strong/anticaching/multi"
	rdb "github.com/WiggidyW/eve-trading-co-go/client/remotedb/internal"
)

type SMAC_WriteShopPurchaseClient = smac.StrongMultiAntiCachingClient[
	WriteShopPurchaseParams,
	time.Time,
	WriteShopPurchaseClient,
]

type WriteShopPurchaseClient struct {
	Inner *rdb.RemoteDBClient
}

// returns the time the appraisal was saved
func (wspc WriteShopPurchaseClient) Fetch(
	ctx context.Context,
	params WriteShopPurchaseParams,
) (*time.Time, error) {
	if err := SaveShopPurchase(
		wspc.Inner,
		ctx,
		params.Appraisal,
	); err != nil {
		return nil, err
	}

	now := time.Now()
	return &now, nil
}
