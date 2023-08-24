package filterregion

import (
	"github.com/WiggidyW/eve-trading-co-go/client/cachekeys"
)

type FilterRegionMarketParams struct {
	RegionId   int32
	TypeId     int32
	IsBuy      bool
	LocationId int64
}

func (p FilterRegionMarketParams) CacheKey() string {
	return cachekeys.FilterRegionMarketCacheKey(
		p.RegionId,
		p.TypeId,
		p.IsBuy,
		p.LocationId,
	)
}
