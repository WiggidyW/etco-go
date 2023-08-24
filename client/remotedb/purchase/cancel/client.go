package cancel

import (
	"context"

	smac "github.com/WiggidyW/eve-trading-co-go/client/caching/strong/anticaching/multi"
	rdb "github.com/WiggidyW/eve-trading-co-go/client/remotedb/internal"
)

type SMAC_CancelShopPurchaseClient = smac.StrongMultiAntiCachingClient[
	CancelShopPurchaseParams,
	struct{},
	CancelShopPurchaseClient,
]

type CancelShopPurchaseClient struct {
	Inner *rdb.RemoteDBClient
}

// returns the time the appraisal was saved
func (wspc CancelShopPurchaseClient) Fetch(
	ctx context.Context,
	params CancelShopPurchaseParams,
) (*struct{}, error) {
	if err := CancelShopPurchase(
		wspc.Inner,
		ctx,
		params.CharacterId,
		params.AppraisalCode,
	); err != nil {
		return nil, err
	}

	return &struct{}{}, nil
}
