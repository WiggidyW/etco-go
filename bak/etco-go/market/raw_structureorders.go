package market

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/esi"
)

const (
	RAW_STRUCTURE_BUF_CAP int = 0
)

var (
	RAW_STRUCTURE_MIN_EXPIRES time.Duration = 30 * time.Minute
)

func init() {
	keys.TypeStrNSStructureMarketOrders = cache.RegisterType[struct{}]("structuremarketorders", 0)
	keys.TypeStrStructureMarketOrders = cache.RegisterType[structureMarketOrdersMap]("structuremarketorders", RAW_STRUCTURE_BUF_CAP)
}

func GetStructureMarketOrders(
	x cache.Context,
	typeId int32,
	isBuy bool,
	structureId int64,
	refreshToken string,
) (
	rep filteredMarketOrders,
	expires time.Time,
	err error,
) {
	nsCacheKey := keys.CacheKeyNSStructureMarketOrders(structureId)
	return rawGet(
		x,
		func(x cache.Context) (
			esi.RepOrStream[esi.OrdersStructureEntry],
			time.Time,
			int,
			error,
		) {
			return esi.GetOrdersStructureEntries(x, structureId, refreshToken)
		},
		newStructureMarketOrdersMap,
		nsCacheKey,
		keys.TypeStrNSStructureMarketOrders,
		keys.CacheKeyStructureMarketOrders(nsCacheKey, typeId, isBuy),
		keys.TypeStrStructureMarketOrders,
		&RAW_STRUCTURE_MIN_EXPIRES,
	)
}

// Data

type typeId = int32
type structureMarketOrdersMap map[typeId]structureOrders
type structureOrders struct {
	Buy  *[]marketOrder
	Sell *[]marketOrder
}

func newStructureMarketOrdersMap() structureMarketOrdersMap {
	return make(structureMarketOrdersMap)
}
func (m structureMarketOrdersMap) GetDiscriminatedOrders(
	entry esi.OrdersStructureEntry,
) *[]marketOrder {
	typeOrders, ok := m[entry.TypeId]
	if !ok {
		typeOrders = structureOrders{
			Buy:  &[]marketOrder{},
			Sell: &[]marketOrder{},
		}
		m[entry.TypeId] = typeOrders
	}
	if entry.IsBuyOrder {
		return typeOrders.Buy
	} else {
		return typeOrders.Sell
	}
}
func (m structureMarketOrdersMap) GetAll(
	nsCacheKey string,
) []marketOrdersWithCacheKey {
	rep := make([]marketOrdersWithCacheKey, 0, len(m)*2)
	for typeId, orders := range m {
		if *orders.Buy != nil {
			rep = append(rep, marketOrdersWithCacheKey{
				CacheKey: keys.CacheKeyStructureMarketOrders(
					nsCacheKey,
					typeId,
					true,
				),
				Orders: orders.Buy,
			})
		}
		if *orders.Sell != nil {
			rep = append(rep, marketOrdersWithCacheKey{
				CacheKey: keys.CacheKeyStructureMarketOrders(
					nsCacheKey,
					typeId,
					false,
				),
				Orders: orders.Sell,
			})
		}
	}
	return rep
}
