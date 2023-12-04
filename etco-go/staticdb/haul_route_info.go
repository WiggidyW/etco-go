package staticdb

import (
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/haulsystemids"
	kvreader_ "github.com/WiggidyW/etco-go/staticdb/kvreaders_"
)

func UNSAFE_GetHaulRoutes() map[b.HaulRouteSystemsKey]b.HaulRouteInfoIndex {
	return kvreader_.KVReaderHaulRoutes.UnsafeGetInner().UnsafeGetInner()
}

func GetHaulRouteSystems(
	key b.HaulRouteSystemsKey,
) (
	startSystemId b.SystemId,
	endSystemId b.SystemId,
) {
	startSystemIndex, endSystemIndex := b.BytesToUint16Pair(key)
	startSystemId = kvreader_.KVReaderSystemIds.UnsafeGet(int(startSystemIndex))
	endSystemId = kvreader_.KVReaderSystemIds.UnsafeGet(int(endSystemIndex))
	return startSystemId, endSystemId
}

type HaulRouteInfo struct {
	MaxVolume      float64
	MinReward      float64
	MaxReward      float64
	TaxRate        float64 // 0-1
	M3Fee          float64
	CollateralRate float64 // 0-1
	RewardStrategy b.HaulRouteRewardStrategy
	typeMap        map[b.TypeId]b.HaulRouteTypePricing
}

func (hri HaulRouteInfo) GetFeePerM3() float64 { return hri.M3Fee }
func (hri HaulRouteInfo) GetTaxRate() float64  { return hri.TaxRate }

type HaulRoutePricingInfo = PricingInfo

func GetHaulRouteInfo(
	systemIds haulsystemids.HaulSystemIds,
) (
	haulRouteInfo *HaulRouteInfo,
) {
	startSystem, exists := kvreader_.KVReaderSystems.Get(systemIds.Start)
	if !exists {
		return nil
	}
	endSystem, exists := kvreader_.KVReaderSystems.Get(systemIds.End)
	if !exists {
		return nil
	}
	infoIndex, exists := kvreader_.KVReaderHaulRoutes.Get(
		b.Uint16PairToBytes(startSystem.Index, endSystem.Index),
	)
	if !exists {
		return nil
	}
	info := kvreader_.KVReaderHaulRouteInfos.UnsafeGet(int(infoIndex))

	return &HaulRouteInfo{
		MaxVolume:      info.MaxVolume().Float64(),
		MinReward:      info.MinReward().Float64(),
		MaxReward:      info.MaxReward().Float64(),
		TaxRate:        float64(info.TaxRate) / 100.0,
		M3Fee:          float64(info.M3Fee),
		CollateralRate: float64(info.CollateralRate) / 100.0,
		RewardStrategy: info.RewardStrategy,
		typeMap: kvreader_.
			KVReaderHaulRouteTypeMaps.
			UnsafeGet(int(info.TypeMapIndex)),
	}
}

func (hri HaulRouteInfo) GetTypePricingInfo(
	typeId b.TypeId,
) (
	haulRouteTypePricingInfo *HaulRoutePricingInfo,
) {
	v, exists := hri.typeMap[typeId]
	if exists {
		haulRouteTypePricingInfo := unsafeGetPricingInfo(v)
		return &haulRouteTypePricingInfo
	} else {
		return nil
	}
}

func (hri HaulRouteInfo) HasTypePricingInfo(typeId b.TypeId) bool {
	_, exists := hri.typeMap[typeId]
	return exists
}
