package market

import (
	"slices"
	"unsafe"
)

type marketOrder struct {
	Price    float64
	Quantity int64
}

func cmpMarketOrders(a, b marketOrder) int {
	diff := a.Price - b.Price
	return *(*int)(unsafe.Pointer(&diff))
}

func sortMarketOrders(orders []marketOrder) {
	slices.SortFunc(orders, cmpMarketOrders)
}

func dedupSortedMarketOrders(
	sortedOrders []marketOrder,
) (
	dedupedOrders []marketOrder,
	totalQuantity int64,
) {
	if len(sortedOrders) == 0 {
		return nil, 0
	}

	dedupedOrders = make([]marketOrder, 0, len(sortedOrders))
	var prevUnique marketOrder = sortedOrders[0]
	totalQuantity = prevUnique.Quantity

	for _, order := range sortedOrders[1:] {
		if order.Price == prevUnique.Price {
			prevUnique.Quantity += order.Quantity
		} else {
			dedupedOrders = append(dedupedOrders, prevUnique)
			prevUnique = order
		}
		totalQuantity += order.Quantity
	}
	dedupedOrders = append(dedupedOrders, prevUnique)

	return dedupedOrders, totalQuantity
}
