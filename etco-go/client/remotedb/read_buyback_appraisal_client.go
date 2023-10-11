package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	wc "github.com/WiggidyW/etco-go/client/caching/weak"
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

const (
	READ_B_APPRAISAL_EXPIRES        time.Duration = 48 * time.Hour
	READ_B_APPRAISAL_MIN_EXPIRES    time.Duration = 0
	READ_B_APPRAISAL_SLOCK_TTL      time.Duration = 30 * time.Second
	READ_B_APPRAISAL_SLOCK_MAX_WAIT time.Duration = 10 * time.Second
)

type WC_ReadBuybackAppraisalClient = wc.WeakCachingClient[
	ReadBuybackAppraisalParams,
	*rdb.BuybackAppraisal,
	cache.ExpirableData[*rdb.BuybackAppraisal],
	ReadBuybackAppraisalClient,
]

func NewWC_ReadBuybackAppraisalClient(
	rdbClient *rdb.RemoteDBClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) WC_ReadBuybackAppraisalClient {
	return wc.NewWeakCachingClient(
		NewReadBuybackAppraisalClient(rdbClient),
		READ_B_APPRAISAL_MIN_EXPIRES,
		cCache,
		sCache,
		READ_B_APPRAISAL_SLOCK_TTL,
		READ_B_APPRAISAL_SLOCK_MAX_WAIT,
	)
}

type ReadBuybackAppraisalClient struct {
	rdbClient *rdb.RemoteDBClient
	expires   time.Duration
}

func NewReadBuybackAppraisalClient(
	rdbClient *rdb.RemoteDBClient,
) ReadBuybackAppraisalClient {
	return ReadBuybackAppraisalClient{
		rdbClient: rdbClient,
		expires:   READ_B_APPRAISAL_EXPIRES,
	}
}

func (rbac ReadBuybackAppraisalClient) Fetch(
	ctx context.Context,
	params ReadBuybackAppraisalParams,
) (*cache.ExpirableData[*rdb.BuybackAppraisal], error) {
	exists, ba, err := rbac.rdbClient.ReadBuybackAppraisal(
		ctx,
		params.AppraisalCode,
	)

	if err != nil {
		return nil, err
	} else if !exists {
		// cache nil so we don't hit the db again
		return cache.NewExpirableDataPtr[*rdb.BuybackAppraisal](
			nil,
			time.Now().Add(rbac.expires),
		), nil
	} else {
		return cache.NewExpirableDataPtr[*rdb.BuybackAppraisal](
			&ba,
			time.Now().Add(rbac.expires),
		), nil
	}
}
