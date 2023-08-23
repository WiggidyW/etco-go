package writebuyback

import (
	"context"
	"time"

	sac "github.com/WiggidyW/weve-esi/client/caching/strong/anticaching/single"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
)

type SAC_WriteBuybackAppraisalClient = sac.StrongAntiCachingClient[
	WriteBuybackAppraisalParams,
	time.Time,
	WriteBuybackAppraisalClient,
]

type WriteBuybackAppraisalClient struct {
	Inner *rdb.RemoteDBClient
}

// returns the time the appraisal was saved
func (wbac WriteBuybackAppraisalClient) Fetch(
	ctx context.Context,
	params WriteBuybackAppraisalParams,
) (*time.Time, error) {
	if err := SaveBuybackAppraisal(
		wbac.Inner,
		ctx,
		params.Appraisal,
	); err != nil {
		return nil, err
	}

	now := time.Now()
	return &now, nil
}
