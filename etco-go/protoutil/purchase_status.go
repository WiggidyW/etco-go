package protoutil

import (
	"github.com/WiggidyW/etco-go/appraisal"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/purchasequeue"
)

func NewPBMakePurchaseStatus(
	rMakePurchaseStatus appraisal.MakePurchaseStatus,
) proto.MakePurchaseStatus {
	switch rMakePurchaseStatus {
	case appraisal.MPS_Success:
		return proto.MakePurchaseStatus_MPS_SUCCESS
	case appraisal.MPS_CooldownLimit:
		return proto.MakePurchaseStatus_MPS_COOLDOWN_LIMIT
	case appraisal.MPS_MaxActiveLimit:
		return proto.MakePurchaseStatus_MPS_MAX_ACTIVE_LIMIT
	case appraisal.MPS_ItemsRejectedAndUnavailable:
		return proto.MakePurchaseStatus_MPS_ITEMS_REJECTED_AND_UNAVAILABLE
	case appraisal.MPS_ItemsRejected:
		return proto.MakePurchaseStatus_MPS_ITEMS_REJECTED
	case appraisal.MPS_ItemsUnavailable:
		return proto.MakePurchaseStatus_MPS_ITEMS_UNAVAILABLE
	}
	panic("unreachable")
}

func NewPBCancelPurchaseStatus(
	rCancelPurchaseStatus purchasequeue.CancelPurchaseStatus,
) proto.CancelPurchaseStatus {
	switch rCancelPurchaseStatus {
	case purchasequeue.CPS_Success:
		return proto.CancelPurchaseStatus_CPS_SUCCESS
	case purchasequeue.CPS_CooldownLimited:
		return proto.CancelPurchaseStatus_CPS_COOLDOWN_LIMIT
	case purchasequeue.CPS_PurchaseNotFound:
		return proto.CancelPurchaseStatus_CPS_NOT_FOUND
	case purchasequeue.CPS_CooldownLimitedAndPurchaseNotFound:
		return proto.CancelPurchaseStatus_CPS_COOLDOWN_LIMIT_AND_NOT_FOUND
	case purchasequeue.CPS_PurchaseNotActive:
		return proto.CancelPurchaseStatus_CPS_NOT_ACTIVE
	}
	panic("unreachable")
}
