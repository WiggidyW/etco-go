package structure

import (
	"github.com/WiggidyW/eve-trading-co-go/client/cachekeys"
	"github.com/WiggidyW/eve-trading-co-go/client/esi/model/ordersstructure"
)

type StructureMarketParams ordersstructure.OrdersStructureParams

func (p StructureMarketParams) CacheKey() string {
	return cachekeys.StructureMarketCacheKey(p.StructureId)
}
