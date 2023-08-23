package shop

import (
	"github.com/WiggidyW/weve-esi/client/appraisal"
	sm "github.com/WiggidyW/weve-esi/client/market/shop"
	"github.com/WiggidyW/weve-esi/staticdb"
)

func NewRejectedAppraisal(params ShopAppraisalParams) *appraisal.ShopAppraisal {
	rItems := params.Items
	aItems := make([]appraisal.ShopItem, 0, len(rItems))
	for _, rItem := range rItems {
		aItems = append(aItems, *sm.NewRejectedItem(
			rItem.TypeId,
			rItem.Quantity,
		))
	}
	return &appraisal.ShopAppraisal{
		// Code: "",
		Items: aItems,
		// Price:       0.0,
		// Time:        time.Time{},
		Version:     staticdb.SHOP_VERSION,
		LocationId:  params.LocationId,
		CharacterId: params.CharacterId,
	}
}
