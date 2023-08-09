package tc

import "github.com/WiggidyW/weve-esi/staticdb"

var KVReaderBuybackInfo tKVReaderBuybackInfo

type tKVReaderBuybackInfo struct{}

func (tKVReaderBuybackInfo) Get(capacity int) *BuybackInfo {
	return newBuybackInfo(capacity)
}

type BuybackInfo struct {
	locationMap map[int32]*BuybackSystemInfo
}

func newBuybackInfo(capacity int) *BuybackInfo {
	return &BuybackInfo{
		locationMap: make(map[int32]*BuybackSystemInfo, capacity),
	}
}

func (bi *BuybackInfo) GetLocation(k int32) (*BuybackSystemInfo, bool) {
	buybackLocation, ok := bi.locationMap[k]
	if !ok {
		if newBuybackLocation, ok := kVReaderBuybackSystem.Get(k); ok {
			bi.locationMap[k] = newBuybackSystemInfo(
				newBuybackLocation,
			)
		} else {
			bi.locationMap[k] = nil
		}
	}
	return buybackLocation, buybackLocation != nil
}

type BuybackSystemInfo struct {
	buybackSystem BuybackSystem
	typeMap       *staticdb.Container[map[int32]BuybackTypeData]
}

func newBuybackSystemInfo(buybackSystem BuybackSystem) *BuybackSystemInfo {
	return &BuybackSystemInfo{buybackSystem: buybackSystem}
}

func (bs *BuybackSystemInfo) M3Fee() float64 {
	return bs.buybackSystem.M3Fee
}

func (bs *BuybackSystemInfo) GetType(k int32) (v *BuybackTypeInfo, ok bool) {
	if bs.typeMap == nil {
		bs.typeMap = staticdb.NewContainer[map[int32]BuybackTypeData](
			kVReaderBuybackTypeMap.UnsafeGet(
				bs.buybackSystem.TypeMapIndex,
			),
		)
	}
	if typeData, ok := bs.typeMap.Inner[k]; ok {
		return &BuybackTypeInfo{
			typeData: typeData,
		}, true
	}
	return nil, false
}

type BuybackTypeInfo struct {
	typeData BuybackTypeData
	pricing  *staticdb.Container[*PricingInfo]
}

func (bt *BuybackTypeInfo) RepEff() (float64, bool) {
	if bt.typeData.ReprocessingEfficiency != 0 {
		return float64(bt.typeData.ReprocessingEfficiency) / 100, true
	}
	return 0, false
}

func (bt *BuybackTypeInfo) RepEffRaw() (uint8, bool) {
	if bt.typeData.ReprocessingEfficiency != 0 {
		return bt.typeData.ReprocessingEfficiency, true
	}
	return 0, false
}

func (bt *BuybackTypeInfo) HasRepEff() bool {
	return bt.typeData.ReprocessingEfficiency != 0
}

func (bt *BuybackTypeInfo) Pricing() (*PricingInfo, bool) {
	if bt.pricing == nil {
		if bt.typeData.PricingIndex != nil {
			p := kVReaderPricing.UnsafeGet(
				*bt.typeData.PricingIndex,
			)
			bt.pricing = staticdb.NewContainer[*PricingInfo](
				newPricingInfo(p),
			)
		} else {
			bt.pricing = staticdb.NewContainer[*PricingInfo](
				nil,
			)
		}
	}
	return bt.pricing.Inner, bt.pricing.Inner != nil
}
