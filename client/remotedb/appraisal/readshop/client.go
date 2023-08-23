package readshop

import (
	"context"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	a "github.com/WiggidyW/weve-esi/client/appraisal"
	wc "github.com/WiggidyW/weve-esi/client/caching/weak"
	rdba "github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
)

type WC_ReadShopAppraisalClient = wc.WeakCachingClient[
	rdba.ReadAppraisalParams,
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
	params rdba.ReadAppraisalParams,
) (*cache.ExpirableData[*a.ShopAppraisal], error) {
	rep := new(a.ShopAppraisal)
	rep.Code = params.AppraisalCode

	if exists, err := rdba.GetAppraisal(
		srac.Inner,
		ctx,
		params.AppraisalCode,
		rdba.SHOP_COLLECTION_ID,
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
