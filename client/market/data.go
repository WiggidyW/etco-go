package market

import (
	"math"

	"github.com/WiggidyW/weve-esi/client/market/internal"
)

type MarketOrder struct {
	Price    float64
	Quantity int64
}

type MarketPrice struct {
	Price float64 // price per 1 item
	Desc  string
}

func UnpackPositivePrice(
	positivePrice *internal.PositivePrice,
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
