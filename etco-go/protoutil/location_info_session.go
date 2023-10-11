package protoutil

import (
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/staticdb"
)

func MaybeNewLocalLocationInfoSession(
	includeLocationInfo bool,
	includeNaming *proto.IncludeLocationNaming,
) *staticdb.LocationInfoSession[*staticdb.LocalLocationNamerTracker] {
	if !includeLocationInfo {
		return nil
	}
	var nameLocation, nameSystem, nameRegion bool
	if includeNaming != nil {
		nameLocation, nameSystem, nameRegion =
			includeNaming.IncludeLocationName,
			includeNaming.IncludeSystemName,
			includeNaming.IncludeRegionName
	} else {
		nameLocation, nameSystem, nameRegion = false, false, false
	}
	namingSessionVal := staticdb.NewLocalLocationInfoSession(
		nameLocation,
		nameSystem,
		nameRegion,
	)
	return &namingSessionVal
}

func MaybeNewSyncLocationInfoSession(
	includeLocationInfo bool,
	includeNaming *proto.IncludeLocationNaming,
) *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker] {
	if !includeLocationInfo {
		return nil
	}
	var nameLocation, nameSystem, nameRegion bool
	if includeNaming != nil {
		nameLocation, nameSystem, nameRegion =
			includeNaming.IncludeLocationName,
			includeNaming.IncludeSystemName,
			includeNaming.IncludeRegionName
	} else {
		nameLocation, nameSystem, nameRegion = false, false, false
	}
	namingSessionVal := staticdb.NewSyncLocationInfoSession(
		nameLocation,
		nameSystem,
		nameRegion,
	)
	return &namingSessionVal
}

func MaybeGetExistingInfoOrTryAddAsStation[LN staticdb.LocationNamerTracker](
	infoSession *staticdb.LocationInfoSession[LN],
	locationId int64,
) (_ *proto.LocationInfo, shouldFetchStructureInfo bool) {
	if infoSession == nil {
		return nil, false
	}
	rLocationInfo := infoSession.GetExistingOrTryAddAsStation(locationId)
	if rLocationInfo != nil {
		return newPBLocationInfo(*rLocationInfo), false
	} else {
		return nil, true
	}
}

func MaybeAddStructureInfo[LN staticdb.LocationNamerTracker](
	infoSession *staticdb.LocationInfoSession[LN],
	locationId int64,
	forbidden bool,
	locationName string,
	systemId int32,
) *proto.LocationInfo {
	if infoSession == nil {
		return nil
	}
	rLocationInfo := infoSession.AddStructure(
		locationId,
		forbidden,
		locationName,
		systemId,
	)
	return newPBLocationInfo(rLocationInfo)
}

func MaybeAddSystem[LN staticdb.LocationNamerTracker](
	infoSession *staticdb.LocationInfoSession[LN],
	systemId int32,
) (regionId int32) {
	if infoSession == nil {
		return 0
	}
	return infoSession.AddSystem(systemId)
}

func MaybeFinishLocationInfoSession[LN staticdb.LocationNamerTracker](
	namingSession *staticdb.LocationInfoSession[LN],
) *proto.LocationNamingMaps {
	if namingSession == nil {
		return nil
	}
	locations, systems, regions := namingSession.Finish()
	return &proto.LocationNamingMaps{
		LocationNames: locations,
		SystemNames:   systems,
		RegionNames:   regions,
	}
}

func newPBLocationInfo(
	rLocationInfo staticdb.LocationInfo,
) *proto.LocationInfo {
	return &proto.LocationInfo{
		IsStructure:        rLocationInfo.IsStructure,
		ForbiddenStructure: rLocationInfo.ForbiddenStructure,
		SystemId:           rLocationInfo.SystemId,
		RegionId:           rLocationInfo.RegionId,
	}
}
