package protoregistry

import (
	"sync"

	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/staticdb"

	b "github.com/WiggidyW/etco-go-bucket"
)

const MaxInt32 int64 = 1<<31 - 1

type ProtoRegistry struct {
	indexMap map[string]int32
	strs     []string
	mu       *sync.RWMutex
}

func NewProtoRegistry(cap int) *ProtoRegistry {
	sr := &ProtoRegistry{
		indexMap: make(map[string]int32, cap),
		strs:     make([]string, 1, cap+1),
		mu:       &sync.RWMutex{},
	}
	sr.strs[0] = "undefined"
	return sr
}

func (r *ProtoRegistry) TryAddStationById(
	locationId int64,
) (
	info *proto.LocationInfo,
) {
	var stationInfo *b.Station = nil
	if locationId <= MaxInt32 {
		stationInfo = staticdb.GetStationInfo(int32(locationId))
	}
	if stationInfo == nil {
		return nil
	} else {
		return &proto.LocationInfo{
			LocationId:         locationId,
			LocationStrIndex:   r.Add(stationInfo.Name),
			IsStructure:        false,
			ForbiddenStructure: false,
			SystemInfo:         r.AddSystemById(stationInfo.SystemId),
		}
	}
}

func (r *ProtoRegistry) AddUndefinedStructure(
	locationId int64,
) (
	info *proto.LocationInfo,
) {
	return &proto.LocationInfo{LocationId: locationId, IsStructure: true}
}

func (r *ProtoRegistry) AddStructure(
	locationId int64,
	name string,
	forbidden bool,
	systemId int32,
) (
	info *proto.LocationInfo,
) {
	return &proto.LocationInfo{
		LocationId:         locationId,
		IsStructure:        true,
		LocationStrIndex:   r.Add(name),
		ForbiddenStructure: forbidden,
		SystemInfo:         r.AddSystemById(systemId),
	}
}

func (r *ProtoRegistry) UNSAFE_AddSystem(
	systemId int32,
	systemData b.System,
) (
	info *proto.SystemInfo,
) {
	info = &proto.SystemInfo{
		SystemId:       systemId,
		SystemStrIndex: r.unsafe_Add(systemData.Name),
		RegionId:       systemData.RegionId,
	}
	regionName := staticdb.GetRegionInfo(systemData.RegionId)
	if regionName != nil {
		info.RegionStrIndex = r.unsafe_Add(*regionName)
	}
	return info
}

func (r *ProtoRegistry) AddUndefinedSystem(
	systemId int32,
) (
	info *proto.SystemInfo,
) {
	return &proto.SystemInfo{SystemId: systemId}
}

func (r *ProtoRegistry) AddSystemById(systemId int32) (
	info *proto.SystemInfo,
) {
	info = &proto.SystemInfo{SystemId: systemId}

	var systemInfo *b.System
	if systemId > 0 {
		systemInfo = staticdb.GetSystemInfo(systemId)
	}

	if systemInfo != nil {
		info.SystemStrIndex = r.Add(systemInfo.Name)
		info.RegionId = systemInfo.RegionId
		regionName := staticdb.GetRegionInfo(systemInfo.RegionId)
		if regionName != nil {
			info.RegionStrIndex = r.Add(*regionName)
		}
	}

	return info
}

func (r *ProtoRegistry) UNSAFE_AddType(
	typeId int32,
	typeData b.TypeData,
) (
	named *proto.NamedTypeId,
) {
	namedType := staticdb.NewSDENamedType(typeData)
	return &proto.NamedTypeId{
		TypeId:                typeId,
		TypeStrIndex:          r.unsafe_Add(namedType.Name),
		CategoryStrIndex:      r.unsafe_Add(namedType.Category),
		GroupStrIndex:         r.unsafe_Add(namedType.Group),
		MarketGroupStrIndexes: r.unsafe_MultiAdd(namedType.MarketGroups),
	}
}

func (r *ProtoRegistry) AddUndefinedType(typeId int32) *proto.NamedTypeId {
	return &proto.NamedTypeId{TypeId: typeId}
}

func (r *ProtoRegistry) AddTypeById(typeId int32) (
	named *proto.NamedTypeId,
) {
	named = &proto.NamedTypeId{TypeId: typeId}

	namedType := staticdb.GetSDENamedType(typeId)
	if namedType != nil {
		named.TypeStrIndex = r.Add(namedType.Name)
		named.CategoryStrIndex = r.Add(namedType.Category)
		named.GroupStrIndex = r.Add(namedType.Group)
		named.MarketGroupStrIndexes = r.MultiAdd(namedType.MarketGroups)
	}

	return named
}

func (r *ProtoRegistry) AddTypeByName(name string) (
	named *proto.NamedTypeId,
	exists bool,
) {
	named = &proto.NamedTypeId{TypeStrIndex: r.Add(name)}

	var namedType *staticdb.SDENamedType
	namedType, named.TypeId = staticdb.GetSDENamedTypeFromName(name)
	if namedType != nil {
		named.CategoryStrIndex = r.Add(namedType.Category)
		named.GroupStrIndex = r.Add(namedType.Group)
		named.MarketGroupStrIndexes = r.MultiAdd(namedType.MarketGroups)
		exists = true
	} else {
		exists = false
	}

	return named, exists
}

func (r *ProtoRegistry) unsafe_MultiAdd(strs []string) (indexes []int32) {
	indexes = make([]int32, len(strs))
	for i, str := range strs {
		indexes[i] = r.unsafe_Add(str)
	}
	return indexes
}

func (r *ProtoRegistry) MultiAdd(strs []string) (indexes []int32) {
	indexes = make([]int32, len(strs))
	for i, str := range strs {
		indexes[i] = r.Add(str)
	}
	return indexes
}

func (r *ProtoRegistry) unsafe_Add(str string) (index int32) {
	var ok bool
	index, ok = r.unsafe_get(str)
	if ok {
		return index
	} else {
		return r.unsafe_addNew(str)
	}
}

func (r *ProtoRegistry) Add(str string) (index int32) {
	var ok bool
	index, ok = r.get(str)
	if ok {
		return index
	} else {
		return r.addNew(str)
	}
}

func (r *ProtoRegistry) get(str string) (index int32, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.unsafe_get(str)
}

func (r *ProtoRegistry) unsafe_get(str string) (index int32, ok bool) {
	index, ok = r.indexMap[str]
	return index, ok
}

func (r *ProtoRegistry) addNew(str string) (index int32) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.unsafe_addNew(str)
}

func (r *ProtoRegistry) unsafe_addNew(str string) (index int32) {
	index = int32(len(r.strs))
	r.indexMap[str] = index
	r.strs = append(r.strs, str)
	return index
}

func (r *ProtoRegistry) Finish() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.strs
}
