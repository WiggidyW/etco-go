package market

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/staticdb"
)

func percentilePriceGet(
	x cache.Context,
	typeId int32,
	pricingInfo staticdb.PricingInfo,
) (
	price float64,
	expires time.Time,
	err error,
) {
	return fetch.HandleFetchVal(
		x,
		nil,
		percentilePriceGetFetchFunc(typeId, pricingInfo),
		nil,
	)
}

func percentilePriceGetFetchFunc(
	typeId int32,
	pricingInfo staticdb.PricingInfo,
) fetch.Fetch[float64] {
	return func(x cache.Context) (
		price *float64,
		expires time.Time,
		_ *postfetch.Params,
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
		if err != nil {
			return nil, expires, nil, err
		}
		price = filteredMarketOrders.percentilePrice(pricingInfo.Percentile)
		return price, expires, nil, nil
	}
}
