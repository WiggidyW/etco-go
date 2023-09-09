package staticdb

import (
	"fmt"
	"sync"

	"github.com/WiggidyW/etco-go/logger"
)

type LocationInfo struct {
	IsStructure        bool
	ForbiddenStructure bool
	SystemId           int32
	RegionId           int32
}

type structureInfo struct {
	Forbidden bool
	SystemId  int32
}

type LocationInfoSession[LN LocationNamerTracker] struct {
	nameLocation         bool
	nameSystem           bool
	nameRegion           bool
	locationNamerTracker LN
}

func NewLocalLocationInfoSession(
	nameLocation bool,
	nameSystem bool,
	nameRegion bool,
) LocationInfoSession[*LocalLocationNamerTracker] {
	return LocationInfoSession[*LocalLocationNamerTracker]{
		nameLocation:         nameLocation,
		nameSystem:           nameSystem,
		nameRegion:           nameRegion,
		locationNamerTracker: newLocalLocationNamerTracker(0, 0, 0, 0, 0),
	}
}

func NewSyncLocationInfoSession(
	nameLocation bool,
	nameSystem bool,
	nameRegion bool,
) LocationInfoSession[*SyncLocationNamerTracker] {
	return LocationInfoSession[*SyncLocationNamerTracker]{
		nameLocation:         nameLocation,
		nameSystem:           nameSystem,
		nameRegion:           nameRegion,
		locationNamerTracker: newSyncLocationNamerTracker(0, 0, 0, 0, 0),
	}
}

func (lis LocationInfoSession[LN]) GetExistingOrTryAddAsStation(
	locationId int64,
) (
	locationInfo *LocationInfo,
) {
	defer logger.Debug(fmt.Sprintf(
		"GetExistingOrTryAddAsStation: %d -> %+v",
		locationId,
		locationInfo,
	))

	var addNames bool

	if lis.locationNamerTracker.hasStation(locationId) {
		// initialize locationInfo with the existing station
		station := *GetStationInfo(int32(locationId))
		locationInfo = &LocationInfo{
			IsStructure:        false,
			ForbiddenStructure: false,
			SystemId:           station.SystemId,
			// RegionId: 0,
		}
		addNames = false

	} else if structure := lis.locationNamerTracker.getStructure(
		locationId,
	); structure != nil {
		// if its an existing unauthorized structure, return now
		if structure.Forbidden {
			return &LocationInfo{
				IsStructure:        true,
				ForbiddenStructure: true,
				SystemId:           -1,
				RegionId:           -1,
			}
		}

		// initialize locationInfo with the existing structure
		locationInfo = &LocationInfo{
			IsStructure:        true,
			ForbiddenStructure: false,
			SystemId:           structure.SystemId,
			// RegionId: 0,
		}
		addNames = false

	} else {
		// if no stations or structures exist, try to get a station
		station := GetStationInfo(int32(locationId))

		// if it's nil, then it's not a station or existing structure
		if station == nil {
			return nil
		}

		// initialize locationInfo with the new station
		locationInfo = &LocationInfo{
			IsStructure:        false,
			ForbiddenStructure: false,
			SystemId:           station.SystemId,
			// RegionId: 0,
		}
		addNames = true

		// add the station's name if we're naming locations
		if lis.nameLocation {
			lis.locationNamerTracker.addLocation(
				locationId,
				station.Name,
			)
		}
	}

	system := *GetSystemInfo(locationInfo.SystemId)
	locationInfo.RegionId = system.RegionId

	if addNames {
		if lis.nameSystem {
			lis.locationNamerTracker.addSystem(
				locationInfo.SystemId,
				system.Name,
			)
		}
		if lis.nameRegion {
			regionName := *GetRegionInfo(system.RegionId)
			lis.locationNamerTracker.addRegion(
				system.RegionId,
				regionName,
			)
		}
	}

	return locationInfo
}

func (lis LocationInfoSession[LN]) AddStructure(
	locationId int64,
	forbidden bool,
	locationName string,
	systemId int32,
) (locationInfo LocationInfo) {
	if structure := lis.locationNamerTracker.getStructure(
		locationId,
	); structure != nil {
		if structure.Forbidden {
			return LocationInfo{
				IsStructure:        true,
				ForbiddenStructure: true,
				SystemId:           -1,
				RegionId:           -1,
			}

		} else {
			system := *GetSystemInfo(structure.SystemId)
			return LocationInfo{
				IsStructure:        true,
				ForbiddenStructure: false,
				SystemId:           structure.SystemId,
				RegionId:           system.RegionId,
			}
		}
	}

	lis.locationNamerTracker.addStructure(locationId, structureInfo{
		Forbidden: forbidden,
		SystemId:  systemId,
	})

	if forbidden {
		return LocationInfo{
			IsStructure:        true,
			ForbiddenStructure: true,
			SystemId:           -1,
			RegionId:           -1,
		}
	}

	if lis.nameLocation {
		lis.locationNamerTracker.addLocation(
			locationId,
			locationName,
		)
	}

	if !lis.nameSystem && !lis.nameRegion {
		return
	}

	system := *GetSystemInfo(systemId)
	if lis.nameSystem {
		lis.locationNamerTracker.addSystem(
			systemId,
			system.Name,
		)
	}
	if lis.nameRegion {
		region := *GetRegionInfo(system.RegionId)
		lis.locationNamerTracker.addRegion(
			system.RegionId,
			region,
		)
	}

	return LocationInfo{
		IsStructure:        true,
		ForbiddenStructure: false,
		SystemId:           systemId,
		RegionId:           system.RegionId,
	}
}

func (lis LocationInfoSession[LN]) AddSystem(systemId int32) (regionId int32) {
	if !lis.nameSystem && !lis.nameRegion {
		return
	}
	system := *GetSystemInfo(systemId)
	if lis.nameSystem {
		lis.locationNamerTracker.addSystem(
			systemId,
			system.Name,
		)
	}
	if lis.nameRegion {
		region := *GetRegionInfo(system.RegionId)
		lis.locationNamerTracker.addRegion(
			system.RegionId,
			region,
		)
	}
	return system.RegionId
}

func (lis LocationInfoSession[LN]) Finish() (
	locationNames map[int64]string,
	systemNames map[int32]string,
	regionNames map[int32]string,
) {
	return lis.locationNamerTracker.finish()
}

type LocationNamerTracker interface {
	hasStation(locationId int64) bool
	addStation(locationId int64)
	getStructure(locationId int64) *structureInfo
	addStructure(locationId int64, structureInfo structureInfo)
	addLocation(locationId int64, name string)
	addSystem(systemId int32, name string)
	addRegion(regionId int32, name string)
	finish() (
		locationNames map[int64]string,
		systemNames map[int32]string,
		regionNames map[int32]string,
	)
}

type LocalLocationNamerTracker struct {
	stations      map[int64]struct{}
	structures    map[int64]structureInfo
	locationNames map[int64]string
	systemNames   map[int32]string
	regionNames   map[int32]string
}

func newLocalLocationNamerTracker(
	capacityStations int,
	capacityStructures int,
	capacityLocationNames int,
	capacitySystemNames int,
	capacityRegionNames int,
) *LocalLocationNamerTracker {
	return &LocalLocationNamerTracker{
		stations: make(map[int64]struct{}, capacityStations),
		structures: make(
			map[int64]structureInfo,
			capacityStructures,
		),
		locationNames: make(map[int64]string, capacityLocationNames),
		systemNames:   make(map[int32]string, capacitySystemNames),
		regionNames:   make(map[int32]string, capacityRegionNames),
	}
}

func (lln *LocalLocationNamerTracker) hasStation(
	locationId int64,
) bool {
	_, ok := lln.stations[locationId]
	return ok
}

func (lln *LocalLocationNamerTracker) addStation(locationId int64) {
	lln.stations[locationId] = struct{}{}
}

func (lln *LocalLocationNamerTracker) getStructure(
	locationId int64,
) *structureInfo {
	if v, ok := lln.structures[locationId]; ok {
		return &v
	} else {
		return nil
	}
}

func (lln *LocalLocationNamerTracker) addStructure(
	locationId int64,
	structureInfo structureInfo,
) {
	lln.structures[locationId] = structureInfo
}

func (lln *LocalLocationNamerTracker) addLocation(
	locationId int64,
	name string,
) {
	lln.locationNames[locationId] = name
}

func (lln *LocalLocationNamerTracker) addSystem(systemId int32, name string) {
	lln.systemNames[systemId] = name
}

func (lln *LocalLocationNamerTracker) addRegion(regionId int32, name string) {
	lln.regionNames[regionId] = name
}

func (lln *LocalLocationNamerTracker) finish() (
	locationNames map[int64]string,
	systemNames map[int32]string,
	regionNames map[int32]string,
) {
	return lln.locationNames, lln.systemNames, lln.regionNames
}

type SyncLocationNamerTracker struct {
	stations            map[int64]struct{}
	rwLockStations      *sync.RWMutex
	structures          map[int64]structureInfo
	rwLockStructures    *sync.RWMutex
	locationNames       map[int64]string
	rwLockLocationNames *sync.RWMutex
	systemNames         map[int32]string
	rwLockSystemNames   *sync.RWMutex
	regionNames         map[int32]string
	rwLockRegionNames   *sync.RWMutex
}

func newSyncLocationNamerTracker(
	capacityStations int,
	capacityStructures int,
	capacityLocationNames int,
	capacitySystemNames int,
	capacityRegionNames int,
) *SyncLocationNamerTracker {
	return &SyncLocationNamerTracker{
		stations: make(
			map[int64]struct{},
			capacityStations,
		),
		rwLockStations: &sync.RWMutex{},
		structures: make(
			map[int64]structureInfo,
			capacityStructures,
		),
		rwLockStructures: &sync.RWMutex{},
		locationNames: make(
			map[int64]string,
			capacityLocationNames,
		),
		rwLockLocationNames: &sync.RWMutex{},
		systemNames: make(
			map[int32]string,
			capacitySystemNames,
		),
		rwLockSystemNames: &sync.RWMutex{},
		regionNames: make(
			map[int32]string,
			capacityRegionNames,
		),
		rwLockRegionNames: &sync.RWMutex{},
	}
}

func (sln *SyncLocationNamerTracker) hasStation(locationId int64) bool {
	sln.rwLockStations.RLock()
	defer sln.rwLockStations.RUnlock()
	_, ok := sln.stations[locationId]
	return ok
}

func (sln *SyncLocationNamerTracker) addStation(locationId int64) {
	sln.rwLockStations.Lock()
	defer sln.rwLockStations.Unlock()
	sln.stations[locationId] = struct{}{}
}

func (sln *SyncLocationNamerTracker) getStructure(
	locationId int64,
) *structureInfo {
	sln.rwLockStructures.RLock()
	defer sln.rwLockStructures.RUnlock()
	if v, ok := sln.structures[locationId]; ok {
		return &v
	} else {
		return nil
	}
}

func (sln *SyncLocationNamerTracker) addStructure(
	locationId int64,
	structureInfo structureInfo,
) {
	sln.rwLockStructures.Lock()
	defer sln.rwLockStructures.Unlock()
	sln.structures[locationId] = structureInfo
}

func (sln *SyncLocationNamerTracker) addLocation(locationId int64, name string) {
	sln.rwLockLocationNames.Lock()
	defer sln.rwLockLocationNames.Unlock()
	sln.locationNames[locationId] = name
}

func (sln *SyncLocationNamerTracker) addSystem(systemId int32, name string) {
	sln.rwLockSystemNames.Lock()
	defer sln.rwLockSystemNames.Unlock()
	sln.systemNames[systemId] = name
}

func (sln *SyncLocationNamerTracker) addRegion(regionId int32, name string) {
	sln.rwLockRegionNames.Lock()
	defer sln.rwLockRegionNames.Unlock()
	sln.regionNames[regionId] = name
}

func (sln *SyncLocationNamerTracker) finish() (
	locationNames map[int64]string,
	systemNames map[int32]string,
	regionNames map[int32]string,
) {
	sln.rwLockLocationNames.RLock()
	defer sln.rwLockLocationNames.RUnlock()
	sln.rwLockSystemNames.RLock()
	defer sln.rwLockSystemNames.RUnlock()
	sln.rwLockRegionNames.RLock()
	defer sln.rwLockRegionNames.RUnlock()
	return sln.locationNames, sln.systemNames, sln.regionNames
}
