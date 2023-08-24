package readbuyback

import (
	"context"
	"time"

	"github.com/WiggidyW/eve-trading-co-go/cache"
	a "github.com/WiggidyW/eve-trading-co-go/client/appraisal"
	wc "github.com/WiggidyW/eve-trading-co-go/client/caching/weak"
	rdba "github.com/WiggidyW/eve-trading-co-go/client/remotedb/appraisal"
	rdb "github.com/WiggidyW/eve-trading-co-go/client/remotedb/internal"
)

type WC_ReadBuybackAppraisalClient = wc.WeakCachingClient[
	rdba.ReadAppraisalParams,
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
	params rdba.ReadAppraisalParams,
) (*cache.ExpirableData[*a.BuybackAppraisal], error) {
	rep := new(a.BuybackAppraisal)
	rep.Code = params.AppraisalCode

	if exists, err := rdba.GetAppraisal(
		brac.Inner,
		ctx,
		params.AppraisalCode,
		rdba.BUYBACK_COLLECTION_ID,
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
