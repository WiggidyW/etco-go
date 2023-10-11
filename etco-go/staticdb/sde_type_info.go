package staticdb

import (
	b "github.com/WiggidyW/etco-go-bucket"
	kvreader_ "github.com/WiggidyW/etco-go/staticdb/kvreaders_"
)

type SDETypeInfo struct {
	ReprocessedMaterials []b.ReprocessedMaterial // maybe nil
	Name                 string                  // english
	Group                string                  // english
	Category             string                  // english
	MarketGroups         []string                // english, top parent last
	Volume               float64
}

func GetSDETypeInfo(typeId b.TypeId) (sdeTypeInfo *SDETypeInfo) {
	typeData, exists := kvreader_.KVReaderTypeDataMap.Get(typeId)
	if !exists {
		return nil
	}
	group := kvreader_.KVReaderGroups.UnsafeGet(typeData.GroupIndex)

	return &SDETypeInfo{
		ReprocessedMaterials: typeData.ReprocessedMaterials,
		Name:                 typeData.Name,
		Group:                group.Name,
		MarketGroups: getMarketGroups(kvreader_.
			KVReaderMarketGroups.
			UnsafeGet(typeData.MarketGroupIndex)),
		Category: kvreader_.
			KVReaderCategories.
			UnsafeGet(group.CategoryIndex),
		Volume: kvreader_.
			KVReaderTypeVolumes.
			UnsafeGet(typeData.VolumeIndex),
	}
}

func getMarketGroups(marketGroup b.MarketGroup) []string {
	marketGroups := make([]string, 0, marketGroup.NumParents+1)
	for marketGroup.ParentIndex != -1 {
		marketGroups = append(marketGroups, marketGroup.Name)
		marketGroup = kvreader_.KVReaderMarketGroups.
			UnsafeGet(marketGroup.ParentIndex)
	}
	marketGroups = append(marketGroups, marketGroup.Name)
	return marketGroups
}

func GetMarketGroupsIndexes(
	marketGroup b.MarketGroup,
	marketGroupIndex int,
) []int32 {
	marketGroupIndexes := make([]int32, 0, marketGroup.NumParents+1)
	marketGroupIndexes = append(
		marketGroupIndexes,
		int32(marketGroupIndex),
	)
	for marketGroup.ParentIndex != -1 {
		marketGroupIndexes = append(
			marketGroupIndexes,
			int32(marketGroup.ParentIndex),
		)
		marketGroup = kvreader_.KVReaderMarketGroups.
			UnsafeGet(marketGroup.ParentIndex)
	}
	return marketGroupIndexes
}
