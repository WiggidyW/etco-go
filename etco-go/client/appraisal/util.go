package appraisal

import (
	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/client/market"
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

func NewRejectedBuybackAppraisal(
	params MakeBuybackAppraisalParams,
) *rdb.BuybackAppraisal {
	rItems := params.Items
	aItems := make([]rdb.BuybackParentItem, 0, len(rItems))
	for _, rItem := range rItems {
		aItems = append(aItems, *market.NewRejectedBuybackItem(
			rItem.TypeId,
			rItem.Quantity,
		))
	}
	return &rdb.BuybackAppraisal{
		// Code: "",
		Items: aItems,
		// Price:       0.0,
		// Fee:         0.0,
		// Time:        time.Time{},
		Version:     build.VERSION_BUYBACK,
		SystemId:    params.SystemId,
		CharacterId: params.CharacterId,
	}
}

func NewRejectedShopAppraisal(
	params MakeShopAppraisalParams,
) *rdb.ShopAppraisal {
	rItems := params.Items
	aItems := make([]rdb.ShopItem, 0, len(rItems))
	for _, rItem := range rItems {
		aItems = append(aItems, *market.NewRejectedShopItem(
			rItem.TypeId,
			rItem.Quantity,
		))
	}
	return &rdb.ShopAppraisal{
		// Code: "",
		Items: aItems,
		// Price:       0.0,
		// Time:        time.Time{},
		Version:     build.VERSION_SHOP,
		LocationId:  params.LocationId,
		CharacterId: params.CharacterId,
	}
}
