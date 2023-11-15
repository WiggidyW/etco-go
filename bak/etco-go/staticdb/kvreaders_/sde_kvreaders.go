package kvreader_

import (
	b "github.com/WiggidyW/etco-go-bucket"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/staticdb/kvreaders_/loader_"
)

var KVReaderCategories loader_.LoadOnceKVReaderGobFSSlice[b.CategoryName] = loader_.
	NewLoadOnceKVReaderGobFSSlice[b.CategoryName](
	b.FILENAME_SDE_CATEGORIES,
	build.CAPACITY_SDE_CATEGORIES,
)

var KVReaderGroups loader_.LoadOnceKVReaderGobFSSlice[b.Group] = loader_.
	NewLoadOnceKVReaderGobFSSlice[b.Group](
	b.FILENAME_SDE_GROUPS,
	build.CAPACITY_SDE_GROUPS,
)

var KVReaderMarketGroups loader_.LoadOnceKVReaderGobFSSlice[b.MarketGroup] = loader_.
	NewLoadOnceKVReaderGobFSSlice[b.MarketGroup](
	b.FILENAME_SDE_MARKET_GROUPS,
	build.CAPACITY_SDE_MARKET_GROUPS,
)

// only includes types with published = true and marketGroupID != null
// includes many languages
var KVReaderNameToTypeId loader_.LoadOnceKVReaderGobFSMap[b.TypeName, b.TypeId] = loader_.
	NewLoadOnceKVReaderGobFSMap[b.TypeName, b.TypeId](
	b.FILENAME_SDE_NAME_TO_TYPE_ID,
	build.CAPACITY_SDE_NAME_TO_TYPE_ID,
)

var KVReaderRegions loader_.LoadOnceKVReaderGobFSMap[b.RegionId, b.RegionName] = loader_.
	NewLoadOnceKVReaderGobFSMap[b.RegionId, b.RegionName](
	b.FILENAME_SDE_REGIONS,
	build.CAPACITY_SDE_REGIONS,
)

var KVReaderSystems loader_.LoadOnceKVReaderGobFSMap[b.SystemId, b.System] = loader_.
	NewLoadOnceKVReaderGobFSMap[b.SystemId, b.System](
	b.FILENAME_SDE_SYSTEMS,
	build.CAPACITY_SDE_SYSTEMS,
)

var KVReaderStations loader_.LoadOnceKVReaderGobFSMap[b.StationId, b.Station] = loader_.
	NewLoadOnceKVReaderGobFSMap[b.StationId, b.Station](
	b.FILENAME_SDE_STATIONS,
	build.CAPACITY_SDE_STATIONS,
)

// only includes types with published = true and marketGroupID != null
var KVReaderTypeDataMap loader_.LoadOnceKVReaderGobFSMap[b.TypeId, b.TypeData] = loader_.
	NewLoadOnceKVReaderGobFSMap[b.TypeId, b.TypeData](
	b.FILENAME_SDE_TYPE_DATA_MAP,
	build.CAPACITY_SDE_TYPE_DATA_MAP,
)

var KVReaderTypeVolumes loader_.LoadOnceKVReaderGobFSSlice[b.TypeVolume] = loader_.
	NewLoadOnceKVReaderGobFSSlice[b.TypeVolume](
	b.FILENAME_SDE_TYPE_VOLUMES,
	build.CAPACITY_SDE_TYPE_VOLUMES,
)

func InitKVReaderCategories() error {
	return KVReaderCategories.Load()
}

func InitKVReaderGroups() error {
	return KVReaderGroups.Load()
}

func InitKVReaderMarketGroups() error {
	return KVReaderMarketGroups.Load()
}

func InitKVReaderNameToTypeId() error {
	return KVReaderNameToTypeId.Load()
}

func InitKVReaderRegions() error {
	return KVReaderRegions.Load()
}

func InitKVReaderSystems() error {
	return KVReaderSystems.Load()
}

func InitKVReaderStations() error {
	return KVReaderStations.Load()
}

func InitKVReaderTypeDataMap() error {
	return KVReaderTypeDataMap.Load()
}

func InitKVReaderTypeVolumes() error {
	return KVReaderTypeVolumes.Load()
}
