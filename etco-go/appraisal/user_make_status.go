package appraisal

import (
	"github.com/WiggidyW/etco-go/proto"
)

type MakePurchaseStatus uint8

const (
	MPS_None MakePurchaseStatus = iota // not attempted / error occured
	MPS_Success
	MPS_CooldownLimit  // appraisal is nil
	MPS_MaxActiveLimit // appraisal is nil
	MPS_ItemsRejectedAndUnavailable
	MPS_ItemsRejected
	MPS_ItemsUnavailable
)

func (mps MakePurchaseStatus) ToProto() proto.MakePurchaseStatus {
	switch mps {
	case MPS_None:
		return proto.MakePurchaseStatus_MPS_NONE
	case MPS_Success:
		return proto.MakePurchaseStatus_MPS_SUCCESS
	case MPS_CooldownLimit:
		return proto.MakePurchaseStatus_MPS_COOLDOWN_LIMIT
	case MPS_MaxActiveLimit:
		return proto.MakePurchaseStatus_MPS_MAX_ACTIVE_LIMIT
	case MPS_ItemsRejectedAndUnavailable:
		return proto.MakePurchaseStatus_MPS_ITEMS_REJECTED_AND_UNAVAILABLE
	case MPS_ItemsRejected:
		return proto.MakePurchaseStatus_MPS_ITEMS_REJECTED
	case MPS_ItemsUnavailable:
		return proto.MakePurchaseStatus_MPS_ITEMS_UNAVAILABLE
	default:
		panic("invalid MakePurchaseStatus")
	}
}
