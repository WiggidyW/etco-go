package staticdb

import (
	b "github.com/WiggidyW/etco-go-bucket"
	kvreader_ "github.com/WiggidyW/etco-go/staticdb/kvreaders_"
)

// mutating the map will possibly cause a panic
func UNSAFE_GetSDETypeDatas() map[b.TypeId]b.TypeData {
	return kvreader_.KVReaderTypeDataMap.UnsafeGetInner().UnsafeGetInner()
}

func GetSDETypeIds() []b.TypeId {
	typeDataMap := UNSAFE_GetSDETypeDatas()
	typeIds := make([]int32, len(typeDataMap))
	i := 0
	for typeId := range typeDataMap {
		typeIds[i] = typeId
		i++
	}
	return typeIds
}

type SDENamedType struct {
	Name         string // english
	Group        string // english
	Category     string // english
	MarketGroups []string
}

func NewSDENamedType(typeData b.TypeData) SDENamedType {
	group := kvreader_.KVReaderGroups.UnsafeGet(typeData.GroupIndex)
	return SDENamedType{
		Name:  typeData.Name,
		Group: group.Name,
		MarketGroups: getMarketGroups(kvreader_.
			KVReaderMarketGroups.
			UnsafeGet(typeData.MarketGroupIndex)),
		Category: kvreader_.
			KVReaderCategories.
			UnsafeGet(group.CategoryIndex),
	}
}

func GetSDENamedTypeFromName(name string) (
	sdeNamedType *SDENamedType,
	typeId b.TypeId,
) {
	var exists bool
	typeId, exists = kvreader_.KVReaderNameToTypeId.Get(name)
	if exists {
		return GetSDENamedType(typeId), typeId
	} else {
		return nil, 0
	}
}

func GetSDENamedType(typeId b.TypeId) *SDENamedType {
	typeData, exists := kvreader_.KVReaderTypeDataMap.Get(typeId)
	if !exists {
		return nil
	} else {
		sdeNamedType := NewSDENamedType(typeData)
		return &sdeNamedType
	}
}

type SDETypeInfo struct {
	ReprocessedMaterials []b.ReprocessedMaterial // maybe nil
	Volume               float64
}

func GetSDETypeInfoFromName(name string) (
	sdeTypeInfo *SDETypeInfo,
	typeId b.TypeId,
) {
	var exists bool
	typeId, exists = kvreader_.KVReaderNameToTypeId.Get(name)
	if exists {
		return GetSDETypeInfo(typeId), typeId
	} else {
		return nil, 0
	}
}

func GetSDETypeInfo(typeId b.TypeId) (sdeTypeInfo *SDETypeInfo) {
	typeData, exists := kvreader_.KVReaderTypeDataMap.Get(typeId)
	if !exists {
		return nil
	}
	return &SDETypeInfo{
		ReprocessedMaterials: typeData.ReprocessedMaterials,
		Volume: kvreader_.KVReaderTypeVolumes.UnsafeGet(
			typeData.VolumeIndex,
		),
	}
}

func getMarketGroups(marketGroup b.MarketGroup) []string {
	marketGroups := make([]string, 0, marketGroup.NumParents+1)
	for marketGroup.ParentIndex != -1 {
		marketGroups = append(marketGroups, marketGroup.Name)
		marketGroup = kvreader_.KVReaderMarketGroups.UnsafeGet(
			marketGroup.ParentIndex,
		)
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
