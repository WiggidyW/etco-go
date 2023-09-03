package raworders_

import (
	"github.com/WiggidyW/etco-go/client/cachekeys"
	"github.com/WiggidyW/etco-go/client/esi/model/ordersstructure"
)

type RawStructureMarketParams ordersstructure.OrdersStructureParams

func (p RawStructureMarketParams) CacheKey() string {
	return cachekeys.StructureMarketCacheKey(p.StructureId)
}
