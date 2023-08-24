package shop

import (
	"github.com/WiggidyW/eve-trading-co-go/client/appraisal"
	sm "github.com/WiggidyW/eve-trading-co-go/client/market/shop"
	"github.com/WiggidyW/eve-trading-co-go/staticdb"
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
