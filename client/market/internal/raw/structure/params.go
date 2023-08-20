package structure

import (
	"github.com/WiggidyW/weve-esi/client/cachekeys"
	"github.com/WiggidyW/weve-esi/client/esi/model/ordersstructure"
)

type StructureMarketParams ordersstructure.OrdersStructureParams

func (p StructureMarketParams) CacheKey() string {
	return cachekeys.StructureMarketCacheKey(p.StructureId)
}
