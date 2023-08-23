package inventory

import (
	"github.com/WiggidyW/weve-esi/client/shopqueue"
	"github.com/WiggidyW/weve-esi/util"
)

type InventoryParams struct {
	LocationId          int64
	ChnSendShopQueueRep *util.ChanSendResult[*shopqueue.ShopQueueResponse]
}
