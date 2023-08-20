package anonymous

import (
	"context"
	"time"

	wb "github.com/WiggidyW/weve-esi/client/remotedb/appraisal/writebuyback"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
)

type WriteBuybackAnonAppraisalClient struct {
	Inner *rdb.RemoteDBClient
}

// returns the time the appraisal was saved
func (wbaac WriteBuybackAnonAppraisalClient) Fetch(
	ctx context.Context,
	params WriteBuybackAnonAppraisalParams,
) (*time.Time, error) {
	if err := wb.SaveBuybackAppraisal(
		wbaac.Inner,
		ctx,
		params.AppraisalCode,
		nil,
		params.Appraisal,
	); err != nil {
		return nil, err
	}

	now := time.Now()
	return &now, nil
}
