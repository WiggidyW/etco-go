package shopqueue

import (
	"github.com/WiggidyW/eve-trading-co-go/client/contracts"
	"github.com/WiggidyW/eve-trading-co-go/util"
)

type ShopQueueResponse struct {
	ParsedShopQueue []string
	Modified        bool // true if the shop queue was modified from its raw state
	ShopContracts   map[string]contracts.Contract
	// shopQueueHashSet map[string]struct{}
}

// func (psq *ParsedShopQueue) ShopQueueHashSet() map[string]struct{} {
// 	if psq.shopQueueHashSet == nil {
// 		psq.shopQueueHashSet = make(
// 			map[string]struct{},
// 			len(psq.ShopQueue),
// 		)
// 		for _, code := range psq.ShopQueue {
// 			psq.shopQueueHashSet[code] = struct{}{}
// 		}
// 	}
// 	return psq.shopQueueHashSet
// }

func (sqr ShopQueueResponse) ShopQueueHashSet() util.MapHashSet[string, struct{}] {
	hs := make(map[string]struct{}, len(sqr.ParsedShopQueue))
	for _, code := range sqr.ParsedShopQueue {
		hs[code] = struct{}{}
	}
	return util.MapHashSet[string, struct{}](hs)
}
