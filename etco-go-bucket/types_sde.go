package etcogobucket

type SDEBucketData struct {
	Categories   []CategoryName
	Groups       []Group
	MarketGroups []MarketGroup
	TypeVolumes  []TypeVolume
	NameToTypeId map[TypeName]TypeId
	Regions      map[RegionId]RegionName
	Systems      map[SystemId]System
	SystemIds    []SystemId
	Stations     map[StationId]Station
	TypeDataMap  map[TypeId]TypeData
	UpdaterData  SDEUpdaterData
}

// type Categories = []CategoryName
// type Groups = []Group
// type MarketGroups = []MarketGroup
// type TypeVolumes = []TypeVolume
// type NameToTypeId = map[TypeName]TypeId
// type Regions = map[RegionId]RegionName
// type Systems = map[SystemId]System
// type Stations = map[StationId]Station
// type TypeDataMap = map[TypeId]TypeData
// type SystemIds = []SystemId

type TypeName = string     // multiple languages
type CategoryName = string // english
type RegionId = int32
type RegionName = string // english
type StationId = int32
type TypeVolume = float64

type Group struct {
	Name          string // english
	CategoryIndex int
}

type MarketGroup struct {
	Name        string // english
	NumParents  uint8
	ParentIndex int // nil -> -1
}

type Station struct {
	SystemId SystemId
	Name     string // english
}

type System struct {
	Index    uint16 // index into SystemIds
	RegionId RegionId
	Name     string // english
}

type TypeData struct {
	ReprocessedMaterials []ReprocessedMaterial
	Name                 string // english
	GroupIndex           int
	MarketGroupIndex     int
	VolumeIndex          int
}

type ReprocessedMaterial struct {
	TypeId   TypeId
	Quantity float64
}
