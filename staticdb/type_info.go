package staticdb

import (
	"github.com/WiggidyW/weve-esi/staticdb/inner/sde"
	"github.com/WiggidyW/weve-esi/staticdb/inner/tc"
)

func GetSDETypeInfo(k int32) *SDETypeInfo {
	if sdeT, ok := sde.KVReaderTypeData.Get(k); ok {
		return newSDETypeInfo(sdeT)
	} else {
		return nil
	}
}

type SDETypeInfo struct {
	ReprMats   []sde.ReprocessedMaterial // maybe nil
	Name       string                    // english
	Group      *string                   // english
	Category   *string                   // english
	MrktGroups []string                  // english // maybe nil // reverse ordered
	Volume     *float64
}

func newSDETypeInfo(sdeT sde.TypeData) *SDETypeInfo {
	typeInfo := &SDETypeInfo{
		ReprMats: sdeT.ReprocessedMaterials,
		Name:     sdeT.Name,
	}

	// set the group if it exists
	if sdeT.GroupIndex != nil {
		group := sde.KVReaderGroups.UnsafeGet(*sdeT.GroupIndex)
		typeInfo.Group = &group.Name

		// set the category if it exists
		if group.CategoryIndex != nil {
			category := sde.KVReaderCategories.UnsafeGet(
				*group.CategoryIndex,
			)
			typeInfo.Category = &category
		}
	}

	// set the market groups if they exist
	if sdeT.MarketGroupIndex != nil {
		// get the first market group
		var mrktGroup sde.MarketGroup
		mrktGroup = sde.KVReaderMarketGroups.UnsafeGet(
			*sdeT.MarketGroupIndex,
		)
		// initialize the slice
		typeInfo.MrktGroups = make([]string, 0, mrktGroup.NumParents+1)
		// append the lowest childs Name
		typeInfo.MrktGroups = append(
			typeInfo.MrktGroups,
			mrktGroup.Name,
		)
		// append all parent Names, top parent last
		for mrktGroup.ParentIndex != nil {
			mrktGroup = sde.KVReaderMarketGroups.UnsafeGet(
				*mrktGroup.ParentIndex,
			)
			typeInfo.MrktGroups = append(
				typeInfo.MrktGroups,
				mrktGroup.Name,
			)
		}
	}

	// set the volume if it exists
	if sdeT.VolumeIndex != nil {
		volume := sde.KVReaderVolumes.UnsafeGet(*sdeT.VolumeIndex)
		typeInfo.Volume = &volume
	}

	return typeInfo
}

func GetBuybackSystemInfo(k int32) *BuybackSystemInfo {
	if tcBS, ok := tc.KVReaderBuybackSystems.Get(k); ok {
		return newBuybackSystemInfo(tcBS)
	} else {
		return nil
	}
}

type BuybackSystemInfo struct {
	M3Fee   *float64
	typeMap map[int32]tc.BuybackTypeData
}

func newBuybackSystemInfo(tcBS tc.BuybackSystem) *BuybackSystemInfo {
	return &BuybackSystemInfo{
		M3Fee:   tcBS.M3Fee,
		typeMap: tc.KVReaderBuybackTypeMap.UnsafeGet(tcBS.TypeMapIndex),
	}
}

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

func GetShopLocationInfo(k int64) *ShopLocationInfo {
	if tcSL, ok := tc.KVReaderShopLocations.Get(k); ok {
		return newShopLocationInfo(tcSL)
	} else {
		return nil
	}
}

type ShopLocationInfo struct {
	BannedFlags HashSet[string] // maybe nil
	typeMap     map[int32]int
}

func newShopLocationInfo(tcSL tc.ShopLocation) *ShopLocationInfo {
	// initialize with the type map
	sli := &ShopLocationInfo{
		typeMap: tc.KVReaderShopTypeMap.UnsafeGet(tcSL.TypeMapIndex),
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

type PricingInfo struct {
	IsBuy            bool
	Prctile          int     // 0 - 100
	Modifier         float64 // 0.01 - 2.55
	MrktName         string
	MrktRefreshToken *string
	MrktLocationId   int64
	MrktIsStructure  bool
}

func newPricingInfo(tcP tc.Pricing) *PricingInfo {
	// validate Prctile
	if tcP.Percentile > 100 {
		panic("Prctile > 100")
	}

	// validate Modifier
	if tcP.Modifier == 0 {
		panic("Modifier == 0")
	}

	// get MarketInfo
	tcM := tc.KVReaderMarket.UnsafeGet(tcP.MarketIndex)

	// validate refresh token
	if tcM.IsStructure && tcM.RefreshToken == nil {
		panic("IsStructure && RefreshToken == nil")
	}

	return &PricingInfo{
		IsBuy:            tcP.IsBuy,
		Prctile:          int(tcP.Percentile),
		Modifier:         float64(tcP.Modifier) / 100,
		MrktName:         tcM.Name,
		MrktRefreshToken: tcM.RefreshToken,
		MrktLocationId:   tcM.LocationId,
		MrktIsStructure:  tcM.IsStructure,
	}
}
