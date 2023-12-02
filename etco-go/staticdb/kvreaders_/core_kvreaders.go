package kvreader_

import (
	b "github.com/WiggidyW/etco-go-bucket"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/staticdb/kvreaders_/loader_"
)

var KVReaderBuybackSystemTypeMaps loader_.LoadOnceKVReaderGobFSSlice[b.BuybackSystemTypeMap] = loader_.
	NewLoadOnceKVReaderGobFSSlice[b.BuybackSystemTypeMap](
	b.FILENAME_CORE_BUYBACK_SYSTEM_TYPE_MAPS,
	build.CAPACITY_CORE_BUYBACK_SYSTEM_TYPE_MAPS,
)

var KVReaderShopLocationTypeMaps loader_.LoadOnceKVReaderGobFSSlice[b.ShopLocationTypeMap] = loader_.
	NewLoadOnceKVReaderGobFSSlice[b.ShopLocationTypeMap](
	b.FILENAME_CORE_SHOP_LOCATION_TYPE_MAPS,
	build.CAPACITY_CORE_SHOP_LOCATION_TYPE_MAPS,
)

var KVReaderHaulRouteTypeMaps loader_.LoadOnceKVReaderGobFSSlice[b.HaulRouteTypeMap] = loader_.
	NewLoadOnceKVReaderGobFSSlice[b.HaulRouteTypeMap](
	b.FILENAME_CORE_HAUL_ROUTE_TYPE_MAPS,
	build.CAPACITY_CORE_HAUL_ROUTE_TYPE_MAPS,
)

var KVReaderBuybackSystems loader_.LoadOnceKVReaderGobFSMap[b.SystemId, b.BuybackSystem] = loader_.
	NewLoadOnceKVReaderGobFSMap[b.SystemId, b.BuybackSystem](
	b.FILENAME_CORE_BUYBACK_SYSTEMS,
	build.CAPACITY_CORE_BUYBACK_SYSTEMS,
)

var KVReaderShopLocations loader_.LoadOnceKVReaderGobFSMap[b.LocationId, b.ShopLocation] = loader_.
	NewLoadOnceKVReaderGobFSMap[b.LocationId, b.ShopLocation](
	b.FILENAME_CORE_SHOP_LOCATIONS,
	build.CAPACITY_CORE_SHOP_LOCATIONS,
)

var KVReaderHaulRoutes loader_.LoadOnceKVReaderGobFSMap[b.HaulRouteSystemsKey, b.HaulRouteInfoIndex] = loader_.
	NewLoadOnceKVReaderGobFSMap[b.HaulRouteSystemsKey, b.HaulRouteInfoIndex](
	b.FILENAME_CORE_HAUL_ROUTES,
	build.CAPACITY_CORE_HAUL_ROUTES,
)

var KVReaderBannedFlagSets loader_.LoadOnceKVReaderGobFSSlice[b.BannedFlagSet] = loader_.
	NewLoadOnceKVReaderGobFSSlice[b.BannedFlagSet](
	b.FILENAME_CORE_BANNED_FLAG_SETS,
	build.CAPACITY_CORE_BANNED_FLAG_SETS,
)

var KVReaderHaulRouteInfos loader_.LoadOnceKVReaderGobFSSlice[b.HaulRouteInfo] = loader_.
	NewLoadOnceKVReaderGobFSSlice[b.HaulRouteInfo](
	b.FILENAME_CORE_HAUL_ROUTE_INFOS,
	build.CAPACITY_CORE_HAUL_ROUTE_INFOS,
)

var KVReaderMarkets loader_.LoadOnceKVReaderGobFSSlice[b.Market] = loader_.
	NewLoadOnceKVReaderGobFSSlice[b.Market](
	b.FILENAME_CORE_MARKETS,
	build.CAPACITY_CORE_MARKETS,
)

var KVReaderPricings loader_.LoadOnceKVReaderGobFSSlice[b.Pricing] = loader_.
	NewLoadOnceKVReaderGobFSSlice[b.Pricing](
	b.FILENAME_CORE_PRICINGS,
	build.CAPACITY_CORE_PRICINGS,
)

func InitKVReaderBuybackSystemTypeMaps() error {
	return KVReaderBuybackSystemTypeMaps.Load()
}

func InitKVReaderShopLocationTypeMaps() error {
	return KVReaderShopLocationTypeMaps.Load()
}

func InitKVReaderBuybackSystems() error {
	return KVReaderBuybackSystems.Load()
}

func InitKVReaderShopLocations() error {
	return KVReaderShopLocations.Load()
}

func InitKVReaderBannedFlagSets() error {
	return KVReaderBannedFlagSets.Load()
}

func InitKVReaderMarkets() error {
	return KVReaderMarkets.Load()
}

func InitKVReaderPricings() error {
	return KVReaderPricings.Load()
}
