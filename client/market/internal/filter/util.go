package filter

import (
	"sort"

	"github.com/WiggidyW/eve-trading-co-go/client/market/internal/raw"
	"github.com/WiggidyW/eve-trading-co-go/util"
)

func SortOrdersByPrice(orders []*raw.MarketOrder) {
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].Price < orders[j].Price
	})
}

// deduplicates and returns the number of nils
func DedupSortedOrders(orders []*raw.MarketOrder) int {
	var nilCount int = 0
	var lastUnique *raw.MarketOrder

	for i := 0; i < len(orders)-1; i++ {
		order := orders[i]
		if lastUnique == nil || order.Price != lastUnique.Price {
			lastUnique = order
		} else {
			lastUnique.Quantity += order.Quantity
			orders[i] = nil
			nilCount++
		}
	}

	return nilCount
}

func SortDedupOrders(orders []raw.MarketOrder) SortedMarketOrders {
	// convert []T to []*T (MUCH faster than sorting []T directly)
	ptrOrders := util.ToPtrSlice(orders)

	// sort the orders by price
	SortOrdersByPrice(ptrOrders)

	// deduplicate the orders with the same price
	nilCount := DedupSortedOrders(ptrOrders)

	// convert []*T back to []T, ignoring nils
	return SortedMarketOrdersFromPtrSlice(ptrOrders, nilCount)
}

// func FromPtrSlice[T any](ptrSlice []*T, nilCount int) []T {
// 	slice := make([]T, 0, len(ptrSlice)-nilCount)

// 	for i := range ptrSlice {
// 		if ptrSlice[i] != nil {
// 			slice = append(slice, *ptrSlice[i])
// 		}
// 	}

// 	return slice
// }
