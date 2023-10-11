package webtocore

import (
	b "github.com/WiggidyW/etco-go-bucket"
)

func convertWebBuybackBuilder(
	webBSTypeMapsBuilder map[b.TypeId]b.WebBuybackSystemTypeBundle,
	corePricings *[]b.Pricing,
	corePricingsIndexMap map[b.Pricing]int,
	coreMarketsIndexMap map[b.MarketName]int,
) (
	coreBSTypeMaps []map[b.TypeId]b.BuybackTypePricing,
	coreBSTypeMapsIndexMap map[b.BundleKey]int,
) {
	return convertWebBuilder(
		webBSTypeMapsBuilder,
		corePricings,
		corePricingsIndexMap,
		coreMarketsIndexMap,
		convertWebBuybackBundle,
	)
}

func convertWebShopBuilder(
	webSLTypeMapsBuilder map[b.TypeId]b.WebShopLocationTypeBundle,
	corePricings *[]b.Pricing,
	corePricingsIndexMap map[b.Pricing]int,
	coreMarketsIndexMap map[b.MarketName]int,
) (
	coreSLTypeMaps []map[b.TypeId]b.ShopTypePricing,
	coreSLTypeMapsIndexMap map[b.BundleKey]int,
) {
	return convertWebBuilder(
		webSLTypeMapsBuilder,
		corePricings,
		corePricingsIndexMap,
		coreMarketsIndexMap,
		convertWebShopBundle,
	)
}

func convertWebBuilder[BP any, V any](
	webTypeMapsBuilder map[b.TypeId]map[b.BundleKey]BP,
	corePricings *[]b.Pricing,
	corePricingsIndexMap map[b.Pricing]int,
	coreMarketsIndexMap map[b.MarketName]int,
	convertBundle func(
		bundlePricing BP,
		corePricings *[]b.Pricing,
		corePricingsIndexMap map[b.Pricing]int,
		coreMarketsIndexMap map[b.MarketName]int,
	) (coreTypeMapValue V),
) (
	coreTypeMaps []map[b.TypeId]V,
	coreTypeMapsIndexMap map[b.BundleKey]int,
) {
	coreTypeMaps, coreTypeMapsIndexMap = initConvertWebBuilder[BP, V](
		webTypeMapsBuilder,
	)
	mutConvertWebBuilder[BP, V](
		webTypeMapsBuilder,
		coreTypeMaps,
		coreTypeMapsIndexMap,
		corePricings,
		corePricingsIndexMap,
		coreMarketsIndexMap,
		convertBundle,
	)
	return coreTypeMaps, coreTypeMapsIndexMap
}

// doesn't perform the full conversion, but rather, simply initializes the maps
// with the correct capacity
func initConvertWebBuilder[BP any, V any](
	webTypeMapsBuilder map[b.TypeId]map[b.BundleKey]BP,
) (
	coreTypeMaps []map[b.TypeId]V, // contains empty maps with correct capacity
	coreTypeMapsIndexMap map[b.BundleKey]int,
) {
	coreTypeMapsIndexMap = make(map[b.BundleKey]int)
	coreTypeMapsCapacities := make([]int, 0) // type map index -> capacity

	for _, bundle := range webTypeMapsBuilder {
		// add the bundle to the index map and increment capacity
		initConvertWebBuilderAddBundle(
			bundle,
			coreTypeMapsIndexMap,
			&coreTypeMapsCapacities,
		)
	}

	// initialize type maps with the correct capacity (= number of unique bundle keys)
	coreTypeMaps = make([]map[b.TypeId]V, 0, len(coreTypeMapsCapacities))

	// initialize the maps with their computed capacities
	for _, capacity := range coreTypeMapsCapacities {
		coreTypeMap := make(map[b.TypeId]V, capacity)
		coreTypeMaps = append(coreTypeMaps, coreTypeMap)
	}

	return coreTypeMaps, coreTypeMapsIndexMap
}

// increment existing capacity /OR/ add a new index entry and set capacity to 1
func initConvertWebBuilderAddBundle[BP any](
	bundle map[b.BundleKey]BP,
	coreTypeMapsIndexMap map[b.BundleKey]int,
	coreTypeMapsCapacities *[]int, // type map index -> type map capacity
) {
	for bundleKey /*, _*/ := range bundle {
		if typeMapIndex, ok := coreTypeMapsIndexMap[bundleKey]; ok {
			// increment the existing capacity
			(*coreTypeMapsCapacities)[typeMapIndex]++
		} else {
			// create a new index entry
			typeMapIndex = len(*coreTypeMapsCapacities)
			coreTypeMapsIndexMap[bundleKey] = typeMapIndex
			// set the capacity to 1
			*coreTypeMapsCapacities = append(
				*coreTypeMapsCapacities,
				1,
			)
		}
	}
}

func mutConvertWebBuilder[BP any, V any](
	webTypeMapsBuilder map[b.TypeId]map[b.BundleKey]BP,
	coreTypeMaps []map[b.TypeId]V, // contains empty maps with correct capacity
	coreTypeMapsIndexMap map[b.BundleKey]int,
	corePricings *[]b.Pricing,
	corePricingsIndexMap map[b.Pricing]int,
	coreMarketsIndexMap map[b.MarketName]int,
	convertBundle func(
		bundlePricing BP,
		corePricings *[]b.Pricing,
		corePricingsIndexMap map[b.Pricing]int,
		coreMarketsIndexMap map[b.MarketName]int,
	) (coreTypeMapValue V),
) {
	for typeId, bundle := range webTypeMapsBuilder {
		for bundleKey, bundlePricing := range bundle {
			coreTypeMapsIndex := coreTypeMapsIndexMap[bundleKey]
			coreTypeMapV := convertBundle(
				bundlePricing,
				corePricings,
				corePricingsIndexMap,
				coreMarketsIndexMap,
			)
			coreTypeMaps[coreTypeMapsIndex][typeId] = coreTypeMapV
		}
	}
}

func convertWebBuybackBundle(
	bundlePricing b.WebBuybackTypePricing,
	corePricings *[]b.Pricing,
	corePricingsIndexMap map[b.Pricing]int,
	coreMarketsIndexMap map[b.MarketName]int,
) (coreTypeMapValue b.BuybackTypePricing) {
	if bundlePricing.Pricing == nil {
		return b.BuybackTypePricing{
			ReprocessingEfficiency: bundlePricing.
				ReprocessingEfficiency,
			PricingIndex: -1,
		}
	}

	corePricing := b.Pricing{
		IsBuy:       bundlePricing.Pricing.IsBuy,
		Percentile:  bundlePricing.Pricing.Percentile,
		Modifier:    bundlePricing.Pricing.Modifier,
		MarketIndex: coreMarketsIndexMap[bundlePricing.Pricing.MarketName],
	}

	var pricingIndex int
	if existingIndex, ok := corePricingsIndexMap[corePricing]; ok {
		pricingIndex = existingIndex
	} else {
		pricingIndex = len(*corePricings)
		corePricingsIndexMap[corePricing] = len(*corePricings)
		*corePricings = append(*corePricings, corePricing)
	}

	return b.BuybackTypePricing{
		ReprocessingEfficiency: bundlePricing.ReprocessingEfficiency,
		PricingIndex:           pricingIndex,
	}
}

func convertWebShopBundle(
	bundlePricing b.WebShopTypePricing,
	corePricings *[]b.Pricing,
	corePricingsIndexMap map[b.Pricing]int,
	coreMarketsIndexMap map[b.MarketName]int,
) (coreTypeMapValue b.ShopTypePricing) {
	corePricing := b.Pricing{
		IsBuy:       bundlePricing.IsBuy,
		Percentile:  bundlePricing.Percentile,
		Modifier:    bundlePricing.Modifier,
		MarketIndex: coreMarketsIndexMap[bundlePricing.MarketName],
	}

	if existingIndex, ok := corePricingsIndexMap[corePricing]; ok {
		coreTypeMapValue = existingIndex
	} else {
		coreTypeMapValue = len(*corePricings)
		corePricingsIndexMap[corePricing] = len(*corePricings)
		*corePricings = append(*corePricings, corePricing)
	}

	return coreTypeMapValue
}
