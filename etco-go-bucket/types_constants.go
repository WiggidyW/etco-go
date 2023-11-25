package etcogobucket

import (
	"time"
)

type ConstantsData struct {
	PURCHASE_MAX_ACTIVE              *int
	MAKE_PURCHASE_COOLDOWN           *time.Duration
	CANCEL_PURCHASE_COOLDOWN         *time.Duration
	CORPORATION_WEB_REFRESH_TOKEN    *string
	STRUCTURE_INFO_WEB_REFRESH_TOKEN *string
	DISCORD_CHANNEL                  *string
	BUYBACK_CONTRACT_NOTIFICATIONS   *bool
	SHOP_CONTRACT_NOTIFICATIONS      *bool
	PURCHASE_NOTIFICATIONS           *bool
}
