package raworders_

import (
	"github.com/WiggidyW/etco-go/client/cachekeys"
	"github.com/WiggidyW/etco-go/client/esi/model/ordersregion"
)

type RawRegionMarketParams ordersregion.OrdersRegionParams

func (p RawRegionMarketParams) CacheKey() string {
	return cachekeys.RegionMarketCacheKey(
		p.RegionId,
		p.TypeId,
		p.IsBuy,
	)
}
