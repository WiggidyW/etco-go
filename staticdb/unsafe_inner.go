package staticdb

import (
	b "github.com/WiggidyW/etco-go-bucket"
	kvreader_ "github.com/WiggidyW/etco-go/staticdb/kvreaders_"
)

// Mutating any of the received data will cause a panic.

func UnsafeGetSDETypeData() map[b.TypeId]b.TypeData {
	return kvreader_.KVReaderTypeDataMap.UnsafeGetInner().UnsafeGetInner()
}

func UnsafeGetSDECategories() []b.CategoryName {
	return kvreader_.KVReaderCategories.UnsafeGetInner().UnsafeGetInner()
}

func UnsafeGetSDEGroups() []b.Group {
	return kvreader_.KVReaderGroups.UnsafeGetInner().UnsafeGetInner()
}

func UnsafeGetSDEMarketGroups() []b.MarketGroup {
	return kvreader_.KVReaderMarketGroups.UnsafeGetInner().UnsafeGetInner()
}

func UnsafeGetCoreBuybackSystems() map[b.SystemId]b.BuybackSystem {
	return kvreader_.KVReaderBuybackSystems.UnsafeGetInner().UnsafeGetInner()
}

func UnsafeGetSDESystems() map[b.SystemId]b.System {
	return kvreader_.KVReaderSystems.UnsafeGetInner().UnsafeGetInner()
}

func UnsafeGetCoreShopLocations() map[b.LocationId]b.ShopLocation {
	return kvreader_.KVReaderShopLocations.UnsafeGetInner().UnsafeGetInner()
}
