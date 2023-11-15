package market

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/esi"
)

const (
	RAW_REGION_BUF_CAP int = 0
)

func init() {
	keys.TypeStrNSRegionMarketOrders = cache.RegisterType[struct{}]("regionmarketorders", 0)
	keys.TypeStrRegionMarketOrders = cache.RegisterType[regionMarketOrdersMap]("regionmarketorders", RAW_REGION_BUF_CAP)
}

func GetRegionMarketOrders(
	x cache.Context,
	typeId int32,
	isBuy bool,
	locationId int64,
	regionId int32,
) (
	rep filteredMarketOrders,
	expires time.Time,
	err error,
) {
	nsCacheKey := keys.CacheKeyNSRegionMarketOrders(regionId, typeId, isBuy)
	return rawGet(
		x,
		func(x cache.Context) (
			esi.RepOrStream[esi.OrdersRegionEntry],
			time.Time,
			int,
			error,
		) {
			return esi.GetOrdersRegionEntries(x, regionId, typeId, isBuy)
		},
		newRegionMarketOrdersMap,
		nsCacheKey,
		keys.TypeStrNSRegionMarketOrders,
		keys.CacheKeyRegionMarketOrders(nsCacheKey, locationId),
		keys.TypeStrRegionMarketOrders,
		nil,
	)
}

// Data

type locationId = int64
type regionMarketOrdersMap map[locationId]*[]marketOrder

func newRegionMarketOrdersMap() regionMarketOrdersMap {
	return make(regionMarketOrdersMap)
}

func (m regionMarketOrdersMap) GetDiscriminatedOrders(
	entry esi.OrdersRegionEntry,
) *[]marketOrder {
	orders, ok := m[entry.LocationId]
	if !ok {
		ordersVal := make([]marketOrder, 0, 1)
		orders = &ordersVal
		m[entry.LocationId] = orders
	}
	return orders
}

func (m regionMarketOrdersMap) GetAll(
	nsCacheKey string,
) []marketOrdersWithCacheKey {
	rep := make([]marketOrdersWithCacheKey, 0, len(m))
	for locationId, orders := range m {
		rep = append(rep, marketOrdersWithCacheKey{
			CacheKey: keys.CacheKeyRegionMarketOrders(nsCacheKey, locationId),
			Orders:   orders,
		})
	}
	return rep
}
