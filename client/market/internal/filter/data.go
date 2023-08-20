package filter

import (
	"github.com/WiggidyW/weve-esi/client/market/internal/raw"
)

type SortedMarketOrders struct {
	MarketOrders []raw.MarketOrder
	Quantity     int64
}

func SortedMarketOrdersFromPtrSlice(
	orders []*raw.MarketOrder,
	nilCount int,
) SortedMarketOrders {
	var quantity int64 = 0
	sortedOrders := make([]raw.MarketOrder, 0, len(orders)-nilCount)

	for _, order := range orders {
		if order != nil {
			sortedOrders = append(sortedOrders, *order)
			quantity += order.Quantity
		}
	}

	return SortedMarketOrders{
		MarketOrders: sortedOrders,
		Quantity:     quantity,
	}
}

func (smo SortedMarketOrders) PercentilePrice(percentile int) float64 {
	if smo.Quantity == 0 {
		// no orders means 0 price
		return 0
	} else if percentile == 100 {
		// percentile 100 = the highest price
		return smo.MarketOrders[len(smo.MarketOrders)-1].Price
	} else if percentile == 0 || len(smo.MarketOrders) == 1 {
		// percentile 0 or only one order = first order price (AKA lowest price)
		return smo.MarketOrders[0].Price
	}

	var price float64
	var currentQuant int64
	stopAtQuant := smo.Quantity * int64(percentile) / 100

	if stopAtQuant > smo.Quantity/2 { // iterate in reverse
		currentQuant = smo.Quantity
		for i := len(smo.MarketOrders) - 1; i >= 0; i-- {
			order := smo.MarketOrders[i]
			currentQuant -= order.Quantity
			if currentQuant <= stopAtQuant {
				price = order.Price
				break
			}
		}

	} else { // iterate normally
		currentQuant = 0
		for i := 0; i < len(smo.MarketOrders); i++ {
			order := smo.MarketOrders[i]
			currentQuant += order.Quantity
			if currentQuant >= stopAtQuant {
				price = order.Price
				break
			}
		}
	}

	return price
}
