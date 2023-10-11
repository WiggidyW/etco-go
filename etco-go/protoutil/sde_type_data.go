package protoutil

import (
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/staticdb"
	kvreader_ "github.com/WiggidyW/etco-go/staticdb/kvreaders_"
)

func NewSDETypes() []int32 {
	UNSAFE_sdeTypeData := staticdb.UnsafeGetSDETypeData()
	typeIds := make([]int32, 0, len(UNSAFE_sdeTypeData))
	for typeId := range UNSAFE_sdeTypeData {
		typeIds = append(typeIds, typeId)
	}
	return typeIds
}

func NewSDENamedTypes() (
	pbNamedTypes []*proto.NamedType,
	typeNamingLists *proto.TypeNamingLists,
) {
	UNSAFE_sdeTypeData := staticdb.UnsafeGetSDETypeData()
	pbNamedTypes = make([]*proto.NamedType, 0, len(UNSAFE_sdeTypeData))

	for typeId, rTypeData := range UNSAFE_sdeTypeData {
		pbNamedTypes = append(
			pbNamedTypes,
			NewSDENamedType(typeId, rTypeData),
		)
	}

	return pbNamedTypes, &proto.TypeNamingLists{
		Groups:       NewSDEGroupNames(),
		Categories:   staticdb.UnsafeGetSDECategories(),
		MarketGroups: NewSDEMarketGroups(),
	}
}

func NewSDENamedType(
	typeId int32,
	rTypeData b.TypeData,
) (
	pbType *proto.NamedType,
) {
	marketGroup := kvreader_.KVReaderMarketGroups.UnsafeGet(
		rTypeData.MarketGroupIndex,
	)
	group := kvreader_.KVReaderGroups.UnsafeGet(rTypeData.GroupIndex)
	return &proto.NamedType{
		TypeId: typeId,
		TypeNamingIndexes: &proto.TypeNamingIndexes{
			Name:          rTypeData.Name,
			GroupIndex:    int32(rTypeData.GroupIndex),
			CategoryIndex: int32(group.CategoryIndex),
			MarketGroupIndexes: staticdb.GetMarketGroupsIndexes(
				marketGroup,
				rTypeData.MarketGroupIndex,
			),
		},
	}
}

func NewSDEGroupNames() []string {
	UNSAFE_sdeGroups := staticdb.UnsafeGetSDEGroups()
	groupNames := make([]string, 0, len(UNSAFE_sdeGroups))
	for _, group := range UNSAFE_sdeGroups {
		groupNames = append(groupNames, group.Name)
	}
	return groupNames
}

func NewSDEMarketGroups() []string {
	UNSAFE_sdeMarketGroups := staticdb.UnsafeGetSDEMarketGroups()
	marketGroupNames := make([]string, 0, len(UNSAFE_sdeMarketGroups))
	for _, marketGroup := range UNSAFE_sdeMarketGroups {
		marketGroupNames = append(marketGroupNames, marketGroup.Name)
	}
	return marketGroupNames
}
