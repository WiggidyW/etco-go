package cancel

type CancelPurchaseStatus uint8

const (
	Success CancelPurchaseStatus = iota
	CooldownLimited
	PurchaseNotFound // not found for the character, not necessarily that it doesn't exist at all
	CooldownLimitedAndPurchaseNotFound
	PurchaseNotActive
)
