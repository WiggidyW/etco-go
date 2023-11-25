package esi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/WiggidyW/etco-go/cache"
)

const (
	ORDERS_REGION_ENTRIES_METHOD   string = http.MethodGet
	ORDERS_REGION_ENTRIES_PER_PAGE int    = 1000
)

type OrdersRegionEntry struct {
	LocationId   int64   `json:"location_id"`
	Price        float64 `json:"price"`
	VolumeRemain int32   `json:"volume_remain"`
	// Duration   int32 `json:"duration"`
	// IsBuyOrder bool  `json:"is_buy_order"`
	// Issued       time.Time `json:"issued"`
	// MinVolume   int32 `json:"min_volume"`
	// OrderId     int64 `json:"order_id"`
	// Range       string `json:"range"`
	// TypeId       int32 `json:"type_id"`
	// VolumeTotal int32 `json:"volume_total"`
}

func (ore OrdersRegionEntry) GetPrice() float64 {
	return ore.Price
}

func (ore OrdersRegionEntry) GetQuantity() int32 {
	return ore.VolumeRemain
}

func ordersRegionEntriesUrl(regionId, typeId int32, isBuy bool) string {
	return fmt.Sprintf(
		"%s/markets/%d/orders/?datasource=%s&type_id=%d&order_type=%s",
		BASE_URL,
		regionId,
		DATASOURCE,
		typeId,
		isBuyToOrderType(isBuy),
	)
}

func isBuyToOrderType(isBuy bool) string {
	if isBuy {
		return "buy"
	}
	return "sell"
}

func GetOrdersRegionEntries(
	x cache.Context,
	regionId, typeId int32,
	isBuy bool,
) (
	repOrStream RepOrStream[OrdersRegionEntry],
	expires time.Time,
	pages int,
	err error,
) {
	return streamGet[OrdersRegionEntry](
		x,
		ordersRegionEntriesUrl(regionId, typeId, isBuy),
		ORDERS_REGION_ENTRIES_METHOD,
		ORDERS_REGION_ENTRIES_PER_PAGE,
		nil,
	)
}
