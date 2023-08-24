package shop

import (
	"github.com/WiggidyW/eve-trading-co-go/client/appraisal"
	"github.com/WiggidyW/eve-trading-co-go/client/authingfwding"
)

type FWD_ShopAppraisalParams = authingfwding.WithAuthFwdableParams[
	ShopAppraisalParams,
	INNERFWD_ShopAppraisalParams,
]

type INNERFWD_ShopAppraisalParams struct {
	Items       []appraisal.BasicItem
	LocationId  int64
	IncludeCode bool
}

func (f INNERFWD_ShopAppraisalParams) ToInnerParams(
	characterId int32,
) ShopAppraisalParams {
	return ShopAppraisalParams{
		Items:       f.Items,
		LocationId:  f.LocationId,
		CharacterId: characterId,
		IncludeCode: f.IncludeCode,
	}
}

type ShopAppraisalParams struct {
	Items       []appraisal.BasicItem
	LocationId  int64
	CharacterId int32 // optional, just puts it into the response
	IncludeCode bool
}
