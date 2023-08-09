package orders

import (
	"sort"
)

var preCalculatePercentiles = [...]int{0, 100}

type MarketOrder struct {
	Price    float64
	Quantity int64
}

type MarketOrders struct {
	Orders      []MarketOrder // pre-sorted by price (ascending AKA cheapest first)
	Quantity    int64         // total quantity of all orders
	Percentiles [101]*float64
}

func NewMarketOrders(orders []MarketOrder) MarketOrders {
	// convert []T to []*T (MUCH faster than sorting []T directly)
	ptrOrders := toPtrSlice(orders)

	// sort the orders by price
	sort.Slice(ptrOrders, func(i, j int) bool {
		return ptrOrders[i].Price < ptrOrders[j].Price
	})

	// deduplicate the orders with the same price
	var lastUnique *MarketOrder
	for i := 0; i < len(ptrOrders)-1; i++ {
		order := ptrOrders[i]
		if lastUnique == nil || order.Price != lastUnique.Price {
			lastUnique = order
		} else {
			lastUnique.Quantity += order.Quantity
			ptrOrders[i] = nil
		}
	}

	// convert []*T back to []T, ignoring nils
	orders = orders[:0]
	for _, order := range ptrOrders {
		if order != nil {
			orders = append(orders, *order)
		}
	}

	// compute the total quantity
	var quantity int64
	for _, order := range orders {
		quantity += order.Quantity
	}

	// create the market orders
	data := &MarketOrders{
		Orders:      orders,
		Quantity:    quantity,
		Percentiles: [101]*float64{},
	}

	// pre-calculate the percentiles and stash them
	for _, p := range preCalculatePercentiles {
		price := data.computePercentile(p)
		data.Percentiles[p] = &price
	}

	return *data
}

func (m MarketOrders) HasOrders() bool {
	return len(m.Orders) > 0
}

func (m MarketOrders) Percentile(i int) (float64, bool) {
	if i < 0 || i > 100 {
		panic("percentile must be between 0 and 100")
	}

	// check if the percentile is already calculated
	cachedPrice := m.Percentiles[i]
	if cachedPrice != nil {
		return *cachedPrice, true
	} else if m.Quantity == 0 {
		return 0, true
	}

	// calculate the percentile
	calcedPrice := m.computePercentile(i)
	return calcedPrice, false
}

func (m MarketOrders) computePercentile(p int) float64 {
	// use shortcuts for 0 and 100
	if p == 0 {
		return m.Orders[0].Price
	} else if p == 100 {
		return m.Orders[len(m.Orders)-1].Price
	}

	var current int64
	stopAt := m.Quantity * int64(p) / 100

	if stopAt > m.Quantity/2 { // iterate in reverse
		current = m.Quantity
		for i := len(m.Orders) - 1; i >= 0; i-- {
			order := m.Orders[i]
			current -= order.Quantity
			if current <= stopAt {
				return order.Price
			}
		}

	} else { // iterate normally
		current = 0
		for i := 0; i < len(m.Orders); i++ {
			order := m.Orders[i]
			current += order.Quantity
			if current >= stopAt {
				return order.Price
			}
		}
	}

	panic("unreachable")
}
