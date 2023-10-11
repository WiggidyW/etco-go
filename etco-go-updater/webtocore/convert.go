package webtocore

import (
	b "github.com/WiggidyW/etco-go-bucket"
)

func convert(
	webBucketData b.WebBucketData,
) (
	coreBucketData b.CoreBucketData,
	err error,
) {
	coreMarkets, coreMarketsIndexMap := convertWebMarkets(
		webBucketData.Markets,
	)

	corePricings := make([]b.Pricing, 0)
	corePricingsIndexMap := make(map[b.Pricing]int)

	coreBSTypeMaps, coreBSTypeMapsIndexMap := convertWebBuybackBuilder(
		webBucketData.BuybackSystemTypeMapsBuilder,
		&corePricings,
		corePricingsIndexMap,
		coreMarketsIndexMap,
	)

	coreSLTypeMaps, coreSLTypeMapsIndexMap := convertWebShopBuilder(
		webBucketData.ShopLocationTypeMapsBuilder,
		&corePricings,
		corePricingsIndexMap,
		coreMarketsIndexMap,
	)

	coreBuybackSystems, err := convertWebBuybackSystems(
		webBucketData.BuybackSystems,
		coreBSTypeMapsIndexMap,
	)
	if err != nil {
		return b.CoreBucketData{}, err
	}

	coreShopLocations, coreBannedFlagSets, err := convertWebShopLocations(
		webBucketData.ShopLocations,
		coreSLTypeMapsIndexMap,
	)
	if err != nil {
		return b.CoreBucketData{}, err
	}

	return b.CoreBucketData{
		BuybackSystemTypeMaps: coreBSTypeMaps,
		ShopLocationTypeMaps:  coreSLTypeMaps,
		BuybackSystems:        coreBuybackSystems,
		ShopLocations:         coreShopLocations,
		BannedFlagSets:        coreBannedFlagSets,
		Pricings:              corePricings,
		Markets:               coreMarkets,
	}, nil
}
