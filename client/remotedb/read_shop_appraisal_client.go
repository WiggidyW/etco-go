package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	wc "github.com/WiggidyW/etco-go/client/caching/weak"
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

const (
	READ_S_APPRAISAL_EXPIRES        time.Duration = 48 * time.Hour
	READ_S_APPRAISAL_MIN_EXPIRES    time.Duration = 0
	READ_S_APPRAISAL_SLOCK_TTL      time.Duration = 30 * time.Second
	READ_S_APPRAISAL_SLOCK_MAX_WAIT time.Duration = 10 * time.Second
)

type WC_ReadShopAppraisalClient = wc.WeakCachingClient[
	ReadShopAppraisalParams,
	*rdb.ShopAppraisal,
	cache.ExpirableData[*rdb.ShopAppraisal],
	ReadShopAppraisalClient,
]

func NewWC_ReadShopAppraisalClient(
	rdbClient *rdb.RemoteDBClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) WC_ReadShopAppraisalClient {
	return wc.NewWeakCachingClient(
		NewReadShopAppraisalClient(rdbClient),
		READ_S_APPRAISAL_MIN_EXPIRES,
		cCache,
		sCache,
		READ_S_APPRAISAL_SLOCK_TTL,
		READ_S_APPRAISAL_SLOCK_MAX_WAIT,
	)
}

type ReadShopAppraisalClient struct {
	rdbClient *rdb.RemoteDBClient
	expires   time.Duration
}

func NewReadShopAppraisalClient(
	rdbClient *rdb.RemoteDBClient,
) ReadShopAppraisalClient {
	return ReadShopAppraisalClient{
		rdbClient: rdbClient,
		expires:   READ_S_APPRAISAL_EXPIRES,
	}
}

func (rsac ReadShopAppraisalClient) Fetch(
	ctx context.Context,
	params ReadShopAppraisalParams,
) (*cache.ExpirableData[*rdb.ShopAppraisal], error) {
	exists, sa, err := rsac.rdbClient.ReadShopAppraisal(
		ctx,
		params.AppraisalCode,
	)

	if err != nil {
		return nil, err
	} else if !exists {
		// cache nil so we don't hit the db again
		return cache.NewExpirableDataPtr[*rdb.ShopAppraisal](
			nil,
			time.Now().Add(rsac.expires),
		), nil
	} else {
		return cache.NewExpirableDataPtr[*rdb.ShopAppraisal](
			&sa,
			time.Now().Add(rsac.expires),
		), nil
	}
}
