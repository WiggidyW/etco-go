package inventory

import (
	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go/client/shopqueue"
)

type InventoryParams struct {
	LocationId          int64
	ChnSendShopQueueRep *chanresult.ChanSendResult[*shopqueue.ShopQueueResponse]
}
