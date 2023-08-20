package market

import (
	"math"
)

func minPriced(price float64, min float64) float64 {
	if price < min {
		return min
	}
	return price
}

func multedByModifier(price, mod float64) float64 {
	return mod * price / 100
}

func roundedToCents(price float64) float64 {
	return math.Round(price*100) / 100
}
