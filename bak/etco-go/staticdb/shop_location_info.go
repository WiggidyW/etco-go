package staticdb

import (
	b "github.com/WiggidyW/etco-go-bucket"
	kvreader_ "github.com/WiggidyW/etco-go/staticdb/kvreaders_"
)

type ShopLocationInfo struct {
	BannedFlagSet b.BannedFlagSet // maybe nil
	TaxRate       float64         // 0-1
	typeMap       map[b.TypeId]b.ShopTypePricing
}

func (sli ShopLocationInfo) GetFeePerM3() float64 { return 0 }
func (sli ShopLocationInfo) GetTaxRate() float64  { return sli.TaxRate }

type ShopPricingInfo = PricingInfo

func GetShopLocationInfo(
	locationId b.LocationId,
) (
	shopLocationInfo *ShopLocationInfo,
) {
	v, exists := kvreader_.KVReaderShopLocations.Get(locationId)
	if exists {
		var bannedFlagSet b.BannedFlagSet = nil
		if v.BannedFlagSetIndex != -1 {
			bannedFlagSet = kvreader_.
				KVReaderBannedFlagSets.
				UnsafeGet(v.BannedFlagSetIndex)
		}
		return &ShopLocationInfo{
			BannedFlagSet: bannedFlagSet,
			TaxRate:       v.TaxRate,
			typeMap: kvreader_.
				KVReaderShopLocationTypeMaps.
				UnsafeGet(v.TypeMapIndex),
		}
	} else {
		return nil
	}
}

func (sli ShopLocationInfo) GetTypePricingInfo(
	typeId b.TypeId,
) (
	shopTypePricingInfo *ShopPricingInfo,
) {
	v, exists := sli.typeMap[typeId]
	if exists {
		shopTypePricingInfo := unsafeGetPricingInfo(v)
		return &shopTypePricingInfo
	} else {
		return nil
	}
}

func (sli ShopLocationInfo) HasTypePricingInfo(typeId b.TypeId) bool {
	_, exists := sli.typeMap[typeId]
	return exists
}

func (sli ShopLocationInfo) HasBannedFlag(flag b.BannedFlag) bool {
	if sli.BannedFlagSet == nil {
		return false
	} else {
		_, ok := sli.BannedFlagSet[flag]
		return ok
	}
}
