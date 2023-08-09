package market

import (
	"fmt"
	"math"
)

func newDesc(
	market string,
	percentile int,
	modifier uint8,
	isBuy bool,
) string {
	var percentileStr string
	if percentile == 0 {
		if isBuy {
			percentileStr = "MinBuy"
		} else {
			percentileStr = "MinSell"
		}
	} else if percentile == 100 {
		if isBuy {
			percentileStr = "MaxBuy"
		} else {
			percentileStr = "MaxSell"
		}
	} else {
		if isBuy {
			percentileStr = fmt.Sprintf(
				"%dth Percentile Buy",
				percentile,
			)
		} else {
			percentileStr = fmt.Sprintf(
				"%dth Percentile Sell",
				percentile,
			)
		}
	}
	return fmt.Sprintf(
		"%s %d%% of %s",
		market,
		modifier,
		percentileStr,
	)
}

func minPriced(price float64, min float64) float64 {
	if price < min {
		return min
	}
	return price
}

func multedByModifier(price float64, mod uint8) float64 {
	return float64(mod) * price / 100
}

func roundedToCents(price float64) float64 {
	return math.Round(price*100) / 100
}
