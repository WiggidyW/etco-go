package purchasequeue

import (
	"github.com/WiggidyW/etco-go/remotedb"
)

// returns a list of codes that are not present in the map
func newPurchaseQueue[T any](
	rawPurchaseQueue remotedb.RawPurchaseQueue,
	shopContracts map[string]T,
) (
	purchaseQueue PurchaseQueue,
	removed []remotedb.CodeAndLocationId,
) {
	purchaseQueue = make(PurchaseQueue, len(rawPurchaseQueue))
	removed = make([]remotedb.CodeAndLocationId, 0)
	var remove bool
	for locationId, codes := range rawPurchaseQueue {
		kept := make([]string, 0, len(codes))
		for _, code := range codes {
			_, remove = shopContracts[code]
			if remove {
				removed = append(
					removed,
					remotedb.NewCodeAndLocationId(code, locationId),
				)
			} else {
				kept = append(kept, code)
			}
		}
		if len(kept) > 0 {
			purchaseQueue[locationId] = kept
		}
	}
	return purchaseQueue, removed
}
