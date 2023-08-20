package readshop

import (
	"context"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	wc "github.com/WiggidyW/weve-esi/client/caching/weak"
	a "github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
)

type WC_ReadShopAppraisalClient = wc.WeakCachingClient[
	a.ReadAppraisalParams,
	*a.ShopAppraisal,
	cache.ExpirableData[*a.ShopAppraisal],
	ReadShopAppraisalClient,
]

type ReadShopAppraisalClient struct {
	Inner   *rdb.RemoteDBClient
	Expires time.Duration
}

func (srac ReadShopAppraisalClient) Fetch(
	ctx context.Context,
	params a.ReadAppraisalParams,
) (*cache.ExpirableData[*a.ShopAppraisal], error) {
	rep := new(a.ShopAppraisal)

	if exists, err := a.GetAppraisal(
		srac.Inner,
		ctx,
		params.AppraisalCode,
		a.SHOP_COLLECTION_ID,
		rep,
	); err != nil {
		return nil, err
	} else if !exists {
		rep = nil // cache nil
	}

	return cache.NewExpirableDataPtr[*a.ShopAppraisal](
		rep,
		time.Now().Add(srac.Expires),
	), nil
}
