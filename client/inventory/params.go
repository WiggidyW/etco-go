package inventory

import (
	"github.com/WiggidyW/eve-trading-co-go/client/shopqueue"
	"github.com/WiggidyW/eve-trading-co-go/util"
)

type InventoryParams struct {
	LocationId          int64
	ChnSendShopQueueRep *util.ChanSendResult[*shopqueue.ShopQueueResponse]
}
