package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	sac "github.com/WiggidyW/etco-go/client/caching/strong/anticaching/single"
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

type SAC_WriteBuybackAppraisalClient = sac.StrongAntiCachingClient[
	WriteBuybackAppraisalParams,
	time.Time,
	WriteBuybackAppraisalClient,
]

func NewSAC_WriteBuybackAppraisalClient(
	rdbClient *rdb.RemoteDBClient,
	readUserDataAntiCache *cache.StrongAntiCache,
) SAC_WriteBuybackAppraisalClient {
	return sac.NewStrongAntiCachingClient(
		NewWriteBuybackAppraisalClient(rdbClient),
		readUserDataAntiCache,
	)
}

type WriteBuybackAppraisalClient struct {
	rdbClient *rdb.RemoteDBClient
}

func NewWriteBuybackAppraisalClient(
	rdbClient *rdb.RemoteDBClient,
) WriteBuybackAppraisalClient {
	return WriteBuybackAppraisalClient{rdbClient}
}

// returns the time the appraisal was saved
func (wbac WriteBuybackAppraisalClient) Fetch(
	ctx context.Context,
	params WriteBuybackAppraisalParams,
) (*time.Time, error) {
	err := wbac.rdbClient.SaveBuybackAppraisal(ctx, params.Appraisal)

	if err != nil {
		return nil, err
	} else {
		now := time.Now()
		return &now, nil
	}
}
