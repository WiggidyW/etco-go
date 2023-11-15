package staticdb

import (
	b "github.com/WiggidyW/etco-go-bucket"
	kvreader_ "github.com/WiggidyW/etco-go/staticdb/kvreaders_"
)

type PricingInfo struct {
	IsBuy              bool
	Percentile         int     // 0 - 100
	Modifier           float64 // 0.01 - 2.55
	MarketName         b.MarketName
	MarketRefreshToken *string
	MarketLocationId   b.LocationId
	MarketIsStructure  bool
}

func unsafeGetPricingInfo(pricingIndex int) PricingInfo {
	v := kvreader_.KVReaderPricings.UnsafeGet(pricingIndex)
	marketInfo := kvreader_.KVReaderMarkets.UnsafeGet(v.MarketIndex)
	return PricingInfo{
		IsBuy:              v.IsBuy,
		Percentile:         int(v.Percentile),
		Modifier:           float64(v.Modifier) / 100,
		MarketName:         marketInfo.Name,
		MarketRefreshToken: marketInfo.RefreshToken,
		MarketLocationId:   marketInfo.LocationId,
		MarketIsStructure:  marketInfo.IsStructure,
	}
}

func (pi PricingInfo) RegionId() (regionId int32, isStation bool) {
	if pi.MarketIsStructure {
		return 0, false
	}
	stationId := int32(pi.MarketLocationId)
	station := kvreader_.KVReaderStations.UnsafeGet(stationId)
	system := kvreader_.KVReaderSystems.UnsafeGet(station.SystemId)
	return system.RegionId, true
}
