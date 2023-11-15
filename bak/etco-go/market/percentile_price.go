package market

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/staticdb"
)

func GetPercentilePrice(
	x cache.Context,
	typeId int32,
	pricingInfo staticdb.PricingInfo,
) (
	price float64,
	expires time.Time,
	err error,
) {
	return percentilePriceGet(x, typeId, pricingInfo)
}
