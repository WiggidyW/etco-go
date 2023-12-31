package staticdb

import (
	b "github.com/WiggidyW/etco-go-bucket"
	kvreader_ "github.com/WiggidyW/etco-go/staticdb/kvreaders_"
)

func GetBuybackSystemIds() []b.SystemId {
	systemsMap :=
		kvreader_.KVReaderBuybackSystems.UnsafeGetInner().UnsafeGetInner()
	systemIds := make([]int32, len(systemsMap))
	i := 0
	for systemId := range systemsMap {
		systemIds[i] = systemId
		i++
	}
	return systemIds
}

type BuybackSystemInfo struct {
	M3Fee   float64
	TaxRate float64 // 0-1
	typeMap map[b.TypeId]b.BuybackTypePricing
}

func (bsi BuybackSystemInfo) GetFeePerM3() float64 { return bsi.M3Fee }
func (bsi BuybackSystemInfo) GetTaxRate() float64  { return bsi.TaxRate }

type BuybackPricingInfo struct {
	ReprocessingEfficiency float64
	PricingInfo            *PricingInfo
}

func GetBuybackSystemInfo(
	systemId b.SystemId,
) (
	buybackSystemInfo *BuybackSystemInfo,
) {
	v, exists := kvreader_.KVReaderBuybackSystems.Get(systemId)
	if exists {
		return &BuybackSystemInfo{
			M3Fee:   v.M3Fee,
			TaxRate: v.TaxRate,
			typeMap: kvreader_.
				KVReaderBuybackSystemTypeMaps.
				UnsafeGet(v.TypeMapIndex),
		}
	} else {
		return nil
	}
}

func (bsi BuybackSystemInfo) GetTypePricingInfo(
	typeId b.TypeId,
) (
	buybackTypePricingInfo *BuybackPricingInfo,
) {
	v, exists := bsi.typeMap[typeId]
	if exists {
		buybackTypePricingInfo = &BuybackPricingInfo{}
		if v.ReprocessingEfficiency > 0 {
			buybackTypePricingInfo.ReprocessingEfficiency =
				float64(v.ReprocessingEfficiency) / 100
		}
		if v.PricingIndex != -1 {
			pricingInfo := unsafeGetPricingInfo(v.PricingIndex)
			buybackTypePricingInfo.PricingInfo = &pricingInfo
		}
		return buybackTypePricingInfo
	} else {
		return nil
	}
}
