package make

import (
	"time"

	"github.com/WiggidyW/weve-esi/client/appraisal"
)

type MakePurchaseParams struct {
	Items       []appraisal.BasicItem
	LocationId  int64
	CharacterId int32
	Cooldown    time.Duration // time to wait before allowing the purchase
	MaxActive   int           // max number of active purchases allowed
}
