package anonymous

import (
	"context"
	"time"

	a "github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
	wb "github.com/WiggidyW/weve-esi/client/remotedb/appraisal/writebuyback"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
)

type WriteBuybackAnonAppraisalClient[
	B a.IBuybackAppraisal[I],
	I a.IBuybackParentItem[CI],
	CI a.IBuybackChildItem,
] struct {
	Inner *rdb.RemoteDBClient
}

// returns the time the appraisal was saved
func (wbaac WriteBuybackAnonAppraisalClient[B, I, CI]) Fetch(
	ctx context.Context,
	params WriteBuybackAnonAppraisalParams[B, I, CI],
) (*time.Time, error) {
	if err := wb.SaveBuybackAppraisal[B, I, CI](
		wbaac.Inner,
		ctx,
		params.AppraisalCode,
		nil,
		params.IAppraisal,
	); err != nil {
		return nil, err
	}

	now := time.Now()
	return &now, nil
}
