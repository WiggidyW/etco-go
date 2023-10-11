package primarysdedata

import (
	b "github.com/WiggidyW/etco-go-bucket"
)

type RawSDEData struct {
	FSDTypeIds       FSDTypeIds
	FSDTypeMaterials FSDTypeMaterials
	FSDMarketGroups  FSDMarketGroups
	FSDGroupIds      FSDGroupIds
	FSDCategoryIds   FSDCategoryIds
	BSDStaStations   BSDStaStations
}

type BSDStaStations = []StaStation
type StaStation struct {
	StationId int32  `yaml:"stationID"`
	Name      string `yaml:"stationName"`
	SystemId  int32  `yaml:"solarSystemID"`
}

type FSDTypeMaterials = map[b.TypeId]TypeMaterials
type TypeMaterials struct {
	Materials []TypeMaterial `yaml:"materials"`
}
type TypeMaterial struct {
	MaterialTypeId b.TypeId `yaml:"materialTypeID"`
	Quantity       float64  `yaml:"quantity"`
}

type FSDGroupIds = map[GroupId]GroupData
type GroupData struct {
	CategoryId CategoryId `yaml:"categoryID"`
	Name       GroupNames `yaml:"name"`
}
type GroupNames struct {
	En string `yaml:"en"`
}

type FSDCategoryIds = map[CategoryId]CategoryData
type CategoryData struct {
	Name CategoryNames `yaml:"name"`
}
type CategoryNames struct {
	En string `yaml:"en"`
}

type FSDMarketGroups = map[MarketGroupId]MarketGroupData
type MarketGroupData struct {
	ParentGroupId *MarketGroupId   `yaml:"parentGroupID"`
	NameId        MarketGroupNames `yaml:"nameID"`
}
type MarketGroupNames struct {
	En string `yaml:"en"`
}

type FSDTypeIds = map[b.TypeId]TypeData
type TypeData struct {
	GroupId       *GroupId       `yaml:"groupID"`
	Published     bool           `yaml:"published"`
	Volume        float64        `yaml:"volume"`
	MarketGroupId *MarketGroupId `yaml:"marketGroupID"`
	Name          TypeNames      `yaml:"name"`
	PortionSize   *float64       `yaml:"portionSize"`
}
type TypeNames struct {
	En string `yaml:"en"`
	De string `yaml:"de"`
	Es string `yaml:"es"`
	Fr string `yaml:"fr"`
	It string `yaml:"it"`
	Ja string `yaml:"ja"`
	Ru string `yaml:"ru"`
	Zh string `yaml:"zh"`
}

type GroupId int32
type CategoryId int32
type MarketGroupId int32
