package protoutil

import (
	"github.com/WiggidyW/etco-go/client/purchase"
	"github.com/WiggidyW/etco-go/proto"
)

func NewPBMakePurchaseStatus(
	rMakePurchaseStatus purchase.MakePurchaseStatus,
) proto.MakePurchaseStatus {
	switch rMakePurchaseStatus {
	case purchase.MPS_Success:
		return proto.MakePurchaseStatus_MPS_SUCCESS
	case purchase.MPS_CooldownLimit:
		return proto.MakePurchaseStatus_MPS_COOLDOWN_LIMIT
	case purchase.MPS_MaxActiveLimit:
		return proto.MakePurchaseStatus_MPS_MAX_ACTIVE_LIMIT
	case purchase.MPS_ItemsRejectedAndUnavailable:
		return proto.MakePurchaseStatus_MPS_ITEMS_REJECTED_AND_UNAVAILABLE
	case purchase.MPS_ItemsRejected:
		return proto.MakePurchaseStatus_MPS_ITEMS_REJECTED
	case purchase.MPS_ItemsUnavailable:
		return proto.MakePurchaseStatus_MPS_ITEMS_UNAVAILABLE
	}
	panic("unreachable")
}

func NewPBCancelPurchaseStatus(
	rCancelPurchaseStatus purchase.CancelPurchaseStatus,
) proto.CancelPurchaseStatus {
	switch rCancelPurchaseStatus {
	case purchase.CPS_Success:
		return proto.CancelPurchaseStatus_CPS_SUCCESS
	case purchase.CPS_CooldownLimited:
		return proto.CancelPurchaseStatus_CPS_COOLDOWN_LIMIT
	case purchase.CPS_PurchaseNotFound:
		return proto.CancelPurchaseStatus_CPS_NOT_FOUND
	case purchase.CPS_CooldownLimitedAndPurchaseNotFound:
		return proto.CancelPurchaseStatus_CPS_COOLDOWN_LIMIT_AND_NOT_FOUND
	case purchase.CPS_PurchaseNotActive:
		return proto.CancelPurchaseStatus_CPS_NOT_ACTIVE
	}
	panic("unreachable")
}
