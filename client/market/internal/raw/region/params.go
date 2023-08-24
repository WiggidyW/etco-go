package region

import (
	"github.com/WiggidyW/eve-trading-co-go/client/cachekeys"
	"github.com/WiggidyW/eve-trading-co-go/client/esi/model/ordersregion"
)

type RegionMarketParams ordersregion.OrdersRegionParams

func (p RegionMarketParams) CacheKey() string {
	return cachekeys.RegionMarketCacheKey(
		p.RegionId,
		p.TypeId,
		p.IsBuy,
	)
}
