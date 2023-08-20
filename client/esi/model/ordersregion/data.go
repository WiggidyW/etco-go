package ordersregion

type OrdersRegionEntry struct {
	// Duration   int32 `json:"duration"`
	// IsBuyOrder bool  `json:"is_buy_order"`
	// Issued       time.Time `json:"issued"`
	LocationId int64 `json:"location_id"`
	// MinVolume   int32 `json:"min_volume"`
	// OrderId     int64 `json:"order_id"`
	Price float64 `json:"price"`
	// Range       string `json:"range"`
	// TypeId       int32 `json:"type_id"`
	VolumeRemain int32 `json:"volume_remain"`
	// VolumeTotal int32 `json:"volume_total"`
}
