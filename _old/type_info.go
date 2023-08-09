package tc

import "github.com/WiggidyW/weve-esi/staticdb"

type KVReaderTypeInfo struct{}

func (KVReaderTypeInfo) Get(k int32) (v *TypeInfo, ok bool) {
	if typeData, ok := kVReaderTypeData.Get(k); ok {
		return newTypeInfo(typeData), true
	}
	return nil, false
}

func (KVReaderTypeInfo) UnsafeGet(k int32) *TypeInfo {
	typeData, _ := kVReaderTypeData.Get(k)
	return newTypeInfo(typeData)
}

type TypeInfo struct {
	typeData       TypeData
	buybackPricing *staticdb.Container[*PricingInfo]
	shopPricing    *staticdb.Container[*PricingInfo]
}

func (t *TypeInfo) BuybackReprocessingEfficiency() (uint8, bool) {
	if t.typeData.BuybackReprocessingEfficiency != 0 {
		return t.typeData.BuybackReprocessingEfficiency, true
	}
	return 0, false
}

func (t *TypeInfo) BuybackPricing() (*PricingInfo, bool) {
	if t.buybackPricing == nil {
		if t.typeData.BuybackPricingIndex != nil {
			p := kVReaderPricing.UnsafeGet(
				*t.typeData.BuybackPricingIndex,
			)
			t.buybackPricing = staticdb.NewContainer[*PricingInfo](
				newPricingInfo(p),
			)
		} else {
			t.buybackPricing = staticdb.NewContainer[*PricingInfo](
				nil,
			)
		}
	}
	return t.buybackPricing.Inner, t.buybackPricing.Inner != nil
}

func (t *TypeInfo) ShopPricing() (*PricingInfo, bool) {
	if t.shopPricing == nil {
		if t.typeData.BuybackPricingIndex != nil {
			p := kVReaderPricing.UnsafeGet(
				*t.typeData.BuybackPricingIndex,
			)
			t.shopPricing = staticdb.NewContainer[*PricingInfo](
				newPricingInfo(p),
			)
		} else {
			t.shopPricing = staticdb.NewContainer[*PricingInfo](
				nil,
			)
		}
	}
	return t.shopPricing.Inner, t.shopPricing.Inner != nil
}

func newTypeInfo(typeData TypeData) *TypeInfo {
	return &TypeInfo{typeData: typeData}
}
