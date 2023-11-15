package protoutil

import (
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/remotedb"
)

func NewPBUserData(rUserData remotedb.UserData) (
	buybackAppraisals []*proto.BuybackAppraisalStatus,
	shopAppraisals []*proto.ShopAppraisalStatus,
	cancelledPurchase int64,
	madePurchase int64,
) {
	buybackAppraisals = make(
		[]*proto.BuybackAppraisalStatus,
		len(rUserData.BuybackAppraisals),
	)
	shopAppraisals = make(
		[]*proto.ShopAppraisalStatus,
		len(rUserData.ShopAppraisals),
	)
	for i, rBuybackAppraisalStatus := range rUserData.BuybackAppraisals {
		buybackAppraisals[i] = NewPBBuybackAppraisalStatus(rBuybackAppraisalStatus)
	}
	for i, rShopAppraisalStatus := range rUserData.ShopAppraisals {
		shopAppraisals[i] = NewPBShopAppraisalStatus(rShopAppraisalStatus)
	}
	cancelledPurchase = rUserData.CancelledPurchase.Unix()
	madePurchase = rUserData.MadePurchase.Unix()
	return buybackAppraisals, shopAppraisals, cancelledPurchase, madePurchase
}

func NewPBBuybackAppraisalStatus(code string) *proto.BuybackAppraisalStatus {
	return &proto.BuybackAppraisalStatus{
		Code: code,
		// Contract: nil,
	}
}

func NewPBShopAppraisalStatus(code string) *proto.ShopAppraisalStatus {
	return &proto.ShopAppraisalStatus{
		Code: code,
		// Contract: nil,
	}
}
