package region

import (
	"github.com/WiggidyW/weve-esi/client/cachekeys"
	"github.com/WiggidyW/weve-esi/client/esi/model/ordersregion"
)

type RegionMarketParams ordersregion.OrdersRegionParams

func (p RegionMarketParams) CacheKey() string {
	return cachekeys.RegionMarketCacheKey(
		p.RegionId,
		p.TypeId,
		p.IsBuy,
	)
}
