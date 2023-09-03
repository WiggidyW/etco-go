package market

import (
	"math"

	"github.com/WiggidyW/etco-go/client/market/marketprice"
)

func UnpackPositivePrice(
	positivePrice *marketprice.PositivePrice,
) (accepted bool, price float64) {
	if positivePrice == nil {
		return false, 0.0
	} else {
		return true, float64(*positivePrice)
	}
}

func RoundedPrice(price float64) float64 {
	if price < 0.01 {
		return 0.01
	} else {
		return math.Round(price*100.0) / 100.0
	}
}
