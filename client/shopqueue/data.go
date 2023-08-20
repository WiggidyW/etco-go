package shopqueue

import "github.com/WiggidyW/weve-esi/client/contracts"

type ParsedShopQueue struct {
	ShopQueue     []string
	ShopContracts map[string]contracts.Contract
	Modified      bool
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

func (psq ParsedShopQueue) ShopQueueHashSet() map[string]struct{} {
	hs := make(map[string]struct{}, len(psq.ShopQueue))
	for _, code := range psq.ShopQueue {
		hs[code] = struct{}{}
	}
	return hs
}
