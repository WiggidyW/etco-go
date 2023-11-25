package purchasequeue

import "github.com/WiggidyW/etco-go/proto"

type CancelPurchaseStatus uint8

const (
	CPS_Success CancelPurchaseStatus = iota
	CPS_CooldownLimited
	CPS_PurchaseNotFound // not found for the character, not necessarily that it doesn't exist at all
	CPS_CooldownLimitedAndPurchaseNotFound
	CPS_PurchaseNotActive
)

func (cps CancelPurchaseStatus) ToProto() proto.CancelPurchaseStatus {
	switch cps {
	case CPS_Success:
		return proto.CancelPurchaseStatus_CPS_SUCCESS
	case CPS_CooldownLimited:
		return proto.CancelPurchaseStatus_CPS_COOLDOWN_LIMIT
	case CPS_PurchaseNotFound:
		return proto.CancelPurchaseStatus_CPS_NOT_FOUND
	case CPS_CooldownLimitedAndPurchaseNotFound:
		return proto.CancelPurchaseStatus_CPS_COOLDOWN_LIMIT_AND_NOT_FOUND
	case CPS_PurchaseNotActive:
		return proto.CancelPurchaseStatus_CPS_NOT_ACTIVE
	default:
		panic("Unknown CancelPurchaseStatus")
	}
}
