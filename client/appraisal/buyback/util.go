package buyback

import (
	"github.com/WiggidyW/weve-esi/client/appraisal"
	bm "github.com/WiggidyW/weve-esi/client/market/buyback"
	"github.com/WiggidyW/weve-esi/staticdb"
)

func NewRejectedAppraisal(params BuybackAppraisalParams) *appraisal.BuybackAppraisal {
	rItems := params.Items
	aItems := make([]appraisal.BuybackParentItem, 0, len(rItems))
	for _, rItem := range rItems {
		aItems = append(aItems, *bm.NewRejectedItem(
			rItem.TypeId,
			rItem.Quantity,
		))
	}
	return &appraisal.BuybackAppraisal{
		// Code: "",
		Items: aItems,
		// Price:       0.0,
		// Time:        time.Time{},
		Version:     staticdb.BUYBACK_VERSION,
		SystemId:    params.SystemId,
		CharacterId: params.CharacterId,
	}
}
