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
	var filteredMarketOrders filteredMarketOrders
	if pricingInfo.MarketIsStructure {
		filteredMarketOrders, expires, err = GetStructureMarketOrders(
			x,
			typeId,
			pricingInfo.IsBuy,
			pricingInfo.MarketLocationId,
			*pricingInfo.MarketRefreshToken,
		)
	} else /* if !pricingInfo.MarketIsStructure */ {
		regionId, _ := pricingInfo.RegionId()
		filteredMarketOrders, expires, err = GetRegionMarketOrders(
			x,
			typeId,
			pricingInfo.IsBuy,
			pricingInfo.MarketLocationId,
			regionId,
		)
	}
	if err == nil {
		price = filteredMarketOrders.percentilePrice(pricingInfo.Percentile)
	}
	return price, expires, err
}
