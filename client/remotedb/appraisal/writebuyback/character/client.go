package anonymous

import (
	"context"
	"time"

	sac "github.com/WiggidyW/weve-esi/client/caching/strong/anticaching/single"
	wb "github.com/WiggidyW/weve-esi/client/remotedb/appraisal/writebuyback"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
)

type SAC_WriteBuybackCharacterAppraisalClient = sac.StrongAntiCachingClient[
	WriteBuybackCharacterAppraisalParams,
	time.Time,
	WriteBuybackCharacterAppraisalClient,
]

type WriteBuybackCharacterAppraisalClient struct {
	Inner *rdb.RemoteDBClient
}

// returns the time the appraisal was saved
func (wbaac WriteBuybackCharacterAppraisalClient) Fetch(
	ctx context.Context,
	params WriteBuybackCharacterAppraisalParams,
) (*time.Time, error) {
	if err := wb.SaveBuybackAppraisal(
		wbaac.Inner,
		ctx,
		params.AppraisalCode,
		&params.CharacterId,
		params.Appraisal,
	); err != nil {
		return nil, err
	}

	now := time.Now()
	return &now, nil
}
