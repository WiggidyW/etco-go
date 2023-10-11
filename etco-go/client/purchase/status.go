package purchase

type CancelPurchaseStatus uint8

const (
	CPS_Success CancelPurchaseStatus = iota
	CPS_CooldownLimited
	CPS_PurchaseNotFound // not found for the character, not necessarily that it doesn't exist at all
	CPS_CooldownLimitedAndPurchaseNotFound
	CPS_PurchaseNotActive
)

type MakePurchaseStatus uint8

const (
	MPS_Success        MakePurchaseStatus = iota
	MPS_CooldownLimit                     // appraisal is nil
	MPS_MaxActiveLimit                    // appraisal is nil
	MPS_ItemsRejectedAndUnavailable
	MPS_ItemsRejected
	MPS_ItemsUnavailable
)
