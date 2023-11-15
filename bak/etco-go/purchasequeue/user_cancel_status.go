package purchasequeue

type CancelPurchaseStatus uint8

const (
	CPS_Success CancelPurchaseStatus = iota
	CPS_CooldownLimited
	CPS_PurchaseNotFound // not found for the character, not necessarily that it doesn't exist at all
	CPS_CooldownLimitedAndPurchaseNotFound
	CPS_PurchaseNotActive
)
