package buyback

import (
	"github.com/WiggidyW/eve-trading-co-go/client/appraisal"
	"github.com/WiggidyW/eve-trading-co-go/client/authingfwding"
)

type FWD_BuybackAppraisalParams = authingfwding.WithAuthFwdableParams[
	BuybackAppraisalParams,
	INNERFWD_BuybackAppraisalParams,
]

type INNERFWD_BuybackAppraisalParams struct {
	Items    []appraisal.BasicItem
	SystemId int32
	Save     bool
}

func (f INNERFWD_BuybackAppraisalParams) ToInnerParams(
	characterId int32,
) BuybackAppraisalParams {
	return BuybackAppraisalParams{
		Items:       f.Items,
		SystemId:    f.SystemId,
		CharacterId: &characterId,
		Save:        f.Save,
	}
}

type BuybackAppraisalParams struct {
	Items       []appraisal.BasicItem
	SystemId    int32
	CharacterId *int32
	Save        bool
}
