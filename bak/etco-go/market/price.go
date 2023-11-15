package market

import (
	"fmt"
	"math"

	"github.com/WiggidyW/etco-go/logger"
	"github.com/WiggidyW/etco-go/staticdb"
)

func UnpackPositivePrice(
	positivePrice *float64,
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

func sdeTypeInfoOrFatal(kind string, typeId int32) *staticdb.SDETypeInfo {
	t := staticdb.GetSDETypeInfo(typeId)
	if t == nil {
		logger.Fatal(fmt.Sprintf(
			"%s valid type %d not found in sde type info",
			kind,
			typeId,
		))
	}
	return t
}
