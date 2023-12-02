package webtocore

import (
	b "github.com/WiggidyW/etco-go-bucket"
)

func convert(
	webBucketData b.WebBucketData,
	webAttrs WebAttrs,
	sdeSystems map[b.SystemId]b.System,
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

	coreHRTypeMaps, coreHRTypeMapsIndexMap := convertWebHaulRouteBuilder(
		webBucketData.HaulRouteTypeMapsBuilder,
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

	coreHaulRoutes, coreHaulRouteInfos, err := convertWebHaulRoutes(
		webBucketData.HaulRoutes,
		coreHRTypeMapsIndexMap,
		sdeSystems,
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
		HaulRouteTypeMaps:     coreHRTypeMaps,
		HaulRoutes:            coreHaulRoutes,
		HaulRouteInfos:        coreHaulRouteInfos,
		UpdaterData: b.CoreUpdaterData{
			CHECKSUM_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER: webAttrs.CHECKSUM_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
			CHECKSUM_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER:  webAttrs.CHECKSUM_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
			CHECKSUM_WEB_BUYBACK_SYSTEMS:                  webAttrs.CHECKSUM_WEB_BUYBACK_SYSTEMS,
			CHECKSUM_WEB_SHOP_LOCATIONS:                   webAttrs.CHECKSUM_WEB_SHOP_LOCATIONS,
			CHECKSUM_WEB_MARKETS:                          webAttrs.CHECKSUM_WEB_MARKETS,
			CHECKSUM_WEB_HAUL_ROUTE_TYPE_MAPS_BUILDER:     webAttrs.CHECKSUM_WEB_HAUL_ROUTE_TYPE_MAPS_BUILDER,
			CHECKSUM_WEB_HAUL_ROUTES:                      webAttrs.CHECKSUM_WEB_HAUL_ROUTES,

			CAPACITY_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER: len(webBucketData.BuybackSystemTypeMapsBuilder),
			CAPACITY_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER:  len(webBucketData.ShopLocationTypeMapsBuilder),
			CAPACITY_WEB_BUYBACK_SYSTEMS:                  len(webBucketData.BuybackSystems),
			CAPACITY_WEB_SHOP_LOCATIONS:                   len(webBucketData.ShopLocations),
			CAPACITY_WEB_MARKETS:                          len(webBucketData.Markets),
			CAPACITY_WEB_HAUL_ROUTE_TYPE_MAPS_BUILDER:     len(webBucketData.HaulRouteTypeMapsBuilder),
			CAPACITY_WEB_HAUL_ROUTES:                      len(webBucketData.HaulRoutes),

			CAPACITY_CORE_BUYBACK_SYSTEM_TYPE_MAPS: len(coreBSTypeMaps),
			CAPACITY_CORE_SHOP_LOCATION_TYPE_MAPS:  len(coreSLTypeMaps),
			CAPACITY_CORE_BUYBACK_SYSTEMS:          len(coreBuybackSystems),
			CAPACITY_CORE_SHOP_LOCATIONS:           len(coreShopLocations),
			CAPACITY_CORE_MARKETS:                  len(coreMarkets),
			CAPACITY_CORE_BANNED_FLAG_SETS:         len(coreBannedFlagSets),
			CAPACITY_CORE_PRICINGS:                 len(corePricings),
			CAPACITY_CORE_HAUL_ROUTE_TYPE_MAPS:     len(coreHRTypeMaps),
			CAPACITY_CORE_HAUL_ROUTES:              len(coreHaulRoutes),
			CAPACITY_CORE_HAUL_ROUTE_INFOS:         len(coreHaulRouteInfos),

			VERSION_BUYBACK: webAttrs.VERSION_STRING,
			VERSION_SHOP:    webAttrs.VERSION_STRING,
		},
	}, nil
}
