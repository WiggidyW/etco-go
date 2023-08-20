package anonymous

import (
	"context"
	"time"

	sac "github.com/WiggidyW/weve-esi/client/caching/strong/anticaching/single"
	a "github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
	wb "github.com/WiggidyW/weve-esi/client/remotedb/appraisal/writebuyback"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
)

type SAC_WriteBuybackCharacterAppraisalClient[
	B a.IBuybackAppraisal[I],
	I a.IBuybackParentItem[CI],
	CI a.IBuybackChildItem,
] struct {
	sac.StrongAntiCachingClient[
		WriteBuybackCharacterAppraisalParams[B, I, CI],
		time.Time,
		WriteBuybackCharacterAppraisalClient[B, I, CI],
	]
}

type WriteBuybackCharacterAppraisalClient[
	B a.IBuybackAppraisal[I],
	I a.IBuybackParentItem[CI],
	CI a.IBuybackChildItem,
] struct {
	Inner *rdb.RemoteDBClient
}

// returns the time the appraisal was saved
func (wbaac WriteBuybackCharacterAppraisalClient[B, I, CI]) Fetch(
	ctx context.Context,
	params WriteBuybackCharacterAppraisalParams[B, I, CI],
) (*time.Time, error) {
	if err := wb.SaveBuybackAppraisal[B, I, CI](
		wbaac.Inner,
		ctx,
		params.AppraisalCode,
		&params.CharacterId,
		params.IAppraisal,
	); err != nil {
		return nil, err
	}

	now := time.Now()
	return &now, nil
}
