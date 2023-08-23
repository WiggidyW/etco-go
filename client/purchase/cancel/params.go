package cancel

import (
	"time"
)

type CancelPurchaseParams struct {
	AppraisalCode string
	CharacterId   int32
	Cooldown      time.Duration // time to wait before allowing character to cancel a purchase
}
