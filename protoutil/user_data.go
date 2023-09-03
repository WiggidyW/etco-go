package protoutil

import (
	"github.com/WiggidyW/etco-go/client/userdata"
	"github.com/WiggidyW/etco-go/proto"
)

func NewPBUserData(rUserData userdata.UserData) (
	buybackAppraisals []*proto.BuybackAppraisalStatus,
	shopAppraisals []*proto.ShopAppraisalStatus,
	cancelledPurchase int64,
	madePurchase int64,
) {
	buybackAppraisals = make(
		[]*proto.BuybackAppraisalStatus,
		0,
		len(rUserData.BuybackAppraisals),
	)
	for _, rBuybackAppraisalStatus := range rUserData.BuybackAppraisals {
		buybackAppraisals = append(
			buybackAppraisals,
			NewPBBuybackAppraisalStatus(rBuybackAppraisalStatus),
		)
	}
	for _, rShopAppraisalStatus := range rUserData.ShopAppraisals {
		shopAppraisals = append(
			shopAppraisals,
			NewPBShopAppraisalStatus(rShopAppraisalStatus),
		)
	}
	cancelledPurchase = rUserData.CancelledPurchase.Unix()
	madePurchase = rUserData.MadePurchase.Unix()
	return buybackAppraisals, shopAppraisals, cancelledPurchase, madePurchase
}

func NewPBBuybackAppraisalStatus(
	rBuybackAppraisalStatus userdata.BuybackAppraisalStatus,
) *proto.BuybackAppraisalStatus {
	if rBuybackAppraisalStatus.Contract == nil {
		return &proto.BuybackAppraisalStatus{
			Code: rBuybackAppraisalStatus.Code,
			// Contract: nil,
		}
	} else {
		return &proto.BuybackAppraisalStatus{
			Code: rBuybackAppraisalStatus.Code,
			Contract: NewPBContract(
				*rBuybackAppraisalStatus.Contract,
			),
		}
	}
}

func NewPBShopAppraisalStatus(
	rShopAppraisalStatus userdata.ShopAppraisalStatus,
) *proto.ShopAppraisalStatus {
	if rShopAppraisalStatus.Contract == nil {
		return &proto.ShopAppraisalStatus{
			Code: rShopAppraisalStatus.Code,
			// Contract:        nil,
			InPurchaseQueue: rShopAppraisalStatus.InPurchaseQueue,
		}
	} else {
		return &proto.ShopAppraisalStatus{
			Code:            rShopAppraisalStatus.Code,
			Contract:        NewPBContract(*rShopAppraisalStatus.Contract),
			InPurchaseQueue: rShopAppraisalStatus.InPurchaseQueue,
		}
	}
}
