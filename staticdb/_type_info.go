package staticdb

import (
	"github.com/WiggidyW/etco-go/staticdb/inner/sde"
	"github.com/WiggidyW/etco-go/staticdb/inner/tc"
)

func (bsi BuybackSystemInfo) GetTypeInfo(k int32) *BuybackTypeInfo {
	if tcBT, ok := bsi.typeMap[k]; ok {
		return newBuybackTypeInfo(tcBT)
	} else {
		return nil
	}
}

type BuybackTypeInfo struct {
	// never both nil
	ReprEff     *float64 // 0.0 - 1.0
	PricingInfo *PricingInfo
}

func newBuybackTypeInfo(tcBT tc.BuybackTypeData) *BuybackTypeInfo {
	typeInfo := BuybackTypeInfo{}

	// validate RepEff and set it if it isn't 0
	if tcBT.ReprocessingEfficiency > 100 {
		panic("ReprocessingEfficiency > 100")
	} else if tcBT.ReprocessingEfficiency > 0 {
		reprEff := float64(tcBT.ReprocessingEfficiency) / 100
		typeInfo.ReprEff = &reprEff
	}

	// set PricingInfo if it isn't nil
	var pricingInfo *PricingInfo
	if tcBT.PricingIndex != nil {
		typeInfo.PricingInfo = newPricingInfo(
			tc.KVReaderPricing.UnsafeGet(*tcBT.PricingIndex),
		)
	}

	// validate that at least one exists
	if typeInfo.ReprEff == nil && pricingInfo == nil {
		panic("ReprocessingEfficiency == nil && PricingInfo == nil")
	}

	return &typeInfo
}

func newShopLocationInfo(tcSL tc.ShopLocation) *ShopLocationInfo {
	// initialize with the type map
	sli := &ShopLocationInfo{
		typeMap: tc.KVReaderShopLocationTypeMap.UnsafeGet(
			tcSL.TypeMapIndex,
		),
	}

	// set banned flags if not nil
	if tcSL.BannedFlagsIndex != nil {
		sli.BannedFlags = HashSet[string](
			tc.KVReaderBannedFlags.UnsafeGet(
				*tcSL.BannedFlagsIndex,
			),
		)
	}

	return sli
}

func (s ShopLocationInfo) HasTypeInfo(k int32) bool {
	_, ok := s.typeMap[k]
	return ok
}

func (s ShopLocationInfo) GetTypeInfo(k int32) *ShopTypeInfo {
	if tcST, ok := s.typeMap[k]; ok {
		return newShopTypeInfo(tcST)
	} else {
		return nil
	}
}

type ShopTypeInfo = PricingInfo

func newShopTypeInfo(tcST int) *ShopTypeInfo {
	return newPricingInfo(
		tc.KVReaderPricing.UnsafeGet(tcST),
	)
}

func (pi PricingInfo) RegionId() (regionId int32, isStation bool) {
	if pi.MrktIsStructure {
		return 0, false
	}
	stationId := int32(pi.MrktLocationId)

	station, ok := sde.KVReaderStations.Get(stationId)
	if !ok {
		panic("!IsStructure && Station not found")
	}

	system, ok := sde.KVReaderSystems.Get(station.SystemId)
	if !ok {
		panic("!IsStructure && System not found")
	}

	return system.RegionId, true
}
