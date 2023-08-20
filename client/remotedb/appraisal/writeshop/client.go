package writeshop

import (
	"context"
	"time"

	smac "github.com/WiggidyW/weve-esi/client/caching/strong/anticaching/multi"
	a "github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
)

type SMAC_WriteShopPurchaseClient[
	S a.IShopAppraisal[I],
	I a.IShopItem,
] struct {
	smac.StrongMultiAntiCachingClient[
		WriteShopPurchaseParams[S, I],
		time.Time,
		WriteShopPurchaseClient[S, I],
	]
}

type WriteShopPurchaseClient[
	S a.IShopAppraisal[I],
	I a.IShopItem,
] struct {
	Inner *rdb.RemoteDBClient
}

// returns the time the appraisal was saved
func (wspc WriteShopPurchaseClient[S, I]) Fetch(
	ctx context.Context,
	params WriteShopPurchaseParams[S, I],
) (*time.Time, error) {
	if err := SaveShopPurchase[S, I](
		wspc.Inner,
		ctx,
		params.AppraisalCode,
		params.CharacterId,
		params.IAppraisal,
	); err != nil {
		return nil, err
	}

	now := time.Now()
	return &now, nil
}
