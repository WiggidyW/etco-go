package esi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/WiggidyW/etco-go/cache"
)

const (
	ORDERS_STRUCTURE_ENTRIES_METHOD   string = http.MethodGet
	ORDERS_STRUCTURE_ENTRIES_PER_PAGE int    = 1000
)

type OrdersStructureEntry struct {
	IsBuyOrder   bool    `json:"is_buy_order"`
	Price        float64 `json:"price"`
	TypeId       int32   `json:"type_id"`
	VolumeRemain int32   `json:"volume_remain"`
	// Duration   int32 `json:"duration"`
	// Issued       time.Time `json:"issued"`
	// LocationId int64 `json:"location_id"`
	// MinVolume   int32 `json:"min_volume"`
	// OrderId     int64 `json:"order_id"`
	// Range       string `json:"range"`
	// VolumeTotal int32 `json:"volume_total"`
}

func (ose OrdersStructureEntry) GetPrice() float64 {
	return ose.Price
}

func (ose OrdersStructureEntry) GetQuantity() int32 {
	return ose.VolumeRemain
}

func ordersStructureEntriesUrl(structureId int64) string {
	return fmt.Sprintf(
		"%s/markets/structures/%d/?datasource=%s",
		BASE_URL,
		structureId,
		DATASOURCE,
	)
}

func GetOrdersStructureEntries(
	x cache.Context,
	structureId int64,
	refreshToken string,
) (
	repOrStream RepOrStream[OrdersStructureEntry],
	expires time.Time,
	pages int,
	err error,
) {
	return streamGet[OrdersStructureEntry](
		x,
		ordersStructureEntriesUrl(structureId),
		ORDERS_STRUCTURE_ENTRIES_METHOD,
		ORDERS_STRUCTURE_ENTRIES_PER_PAGE,
		EsiAuthMarkets(refreshToken),
	)
}
