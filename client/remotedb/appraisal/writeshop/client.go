package writeshop

import (
	"context"
	"time"

	smac "github.com/WiggidyW/weve-esi/client/caching/strong/anticaching/multi"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
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
		params.AppraisalCode,
		params.CharacterId,
		params.Appraisal,
	); err != nil {
		return nil, err
	}

	now := time.Now()
	return &now, nil
}
