package market

import (
	"github.com/WiggidyW/etco-go/logger"
)

type marketOrdersWithCacheKey struct {
	CacheKey string
	Orders   *[]marketOrder
}

type filteredOrdersWithCacheKey struct {
	CacheKey string
	Orders   *filteredMarketOrders
}

type marketOrdersMap[E any] interface {
	GetDiscriminatedOrders(entry E) *[]marketOrder
	GetAll(nsCacheKey string) []marketOrdersWithCacheKey
}

type marketOrdersEntry interface {
	GetPrice() float64
	GetQuantity() int32
}

func marketOrdersAppendEntry[E marketOrdersEntry](
	marketOrders *[]marketOrder,
	entry E,
) {
	*marketOrders = append(
		*marketOrders,
		marketOrder{
			Price:    entry.GetPrice(),
			Quantity: int64(entry.GetQuantity()),
		},
	)
}

func marketOrdersMapInsertEntries[
	E marketOrdersEntry,
	M marketOrdersMap[E],
](
	marketOrdersMap M,
	entries []E,
) {
	for _, entry := range entries {
		orders := marketOrdersMap.GetDiscriminatedOrders(entry)
		marketOrdersAppendEntry(orders, entry)
	}
}

// filtered

var zeroF64 float64 = 0.0

type filteredMarketOrders struct {
	Orders   []marketOrder
	Quantity int64
}

func newFilteredMarketOrders(rawOrders []marketOrder) *filteredMarketOrders {
	sortMarketOrders(rawOrders)
	orders, quantity := dedupSortedMarketOrders(rawOrders)
	return &filteredMarketOrders{Orders: orders, Quantity: quantity}
}

func (fmo filteredMarketOrders) percentilePrice(
	percentile int,
) (pricePtr *float64) {
	if fmo.Quantity == 0 {
		pricePtr = &zeroF64
	} else if percentile == 100 {
		pricePtr = &fmo.Orders[len(fmo.Orders)-1].Price
	} else if percentile == 0 {
		pricePtr = &fmo.Orders[0].Price
	} else if percentile > 100 || percentile < 0 {
		panic("percentile must be [0 to 100]")
	} else {
		var price float64
		var currentQuantity int64
		stopAtQuant := fmo.Quantity * int64(percentile) / 100

		if stopAtQuant > fmo.Quantity/2 { // iterate in reverse
			currentQuantity = fmo.Quantity
			for i := len(fmo.Orders) - 1; i >= 0; i-- {
				order := fmo.Orders[i]
				currentQuantity -= order.Quantity
				if currentQuantity <= stopAtQuant {
					price = order.Price
					break
				}
			}

		} else { // iterate normally
			currentQuantity = 0
			for i := 0; i < len(fmo.Orders); i++ {
				order := fmo.Orders[i]
				currentQuantity += order.Quantity
				if currentQuantity >= stopAtQuant {
					price = order.Price
					break
				}
			}
		}
		pricePtr = &price
	}
	if *pricePtr < 0 {
		pricePtr = &zeroF64
		logger.Err("Negative price found in percentilePrice")
	}
	return pricePtr
}
