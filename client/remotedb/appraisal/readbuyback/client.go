package readbuyback

import (
	"context"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	wc "github.com/WiggidyW/weve-esi/client/caching/weak"
	a "github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
)

type WC_ReadBuybackAppraisalClient = wc.WeakCachingClient[
	a.ReadAppraisalParams,
	*a.BuybackAppraisal,
	cache.ExpirableData[*a.BuybackAppraisal],
	ReadBuybackAppraisalClient,
]

type ReadBuybackAppraisalClient struct {
	Inner   *rdb.RemoteDBClient
	Expires time.Duration
}

func (brac ReadBuybackAppraisalClient) Fetch(
	ctx context.Context,
	params a.ReadAppraisalParams,
) (*cache.ExpirableData[*a.BuybackAppraisal], error) {
	rep := new(a.BuybackAppraisal)

	if exists, err := a.GetAppraisal(
		brac.Inner,
		ctx,
		params.AppraisalCode,
		a.BUYBACK_COLLECTION_ID,
		rep,
	); err != nil {
		return nil, err
	} else if !exists {
		rep = nil // cache nil
	}

	return cache.NewExpirableDataPtr[*a.BuybackAppraisal](
		rep,
		time.Now().Add(brac.Expires),
	), nil
}
