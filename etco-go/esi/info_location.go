package esi

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"
)

func ProtoGetLocationInfo(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	locationId int64, // station or structure (unknown which)
) (
	info *proto.LocationInfo,
	expires time.Time,
	err error,
) {
	// try to get it as a station
	// stations are locally-present static data
	info = r.TryAddStationById(locationId)
	if info != nil {
		expires = fetch.MAX_EXPIRES
	} else {
		info, expires, err = protoGetStructureInfo(x, r, locationId)
	}
	return info, expires, err
}

func ProtoGetLocationInfoCOV(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	locationId int64, // station or structure (unknown which)
) (
	infoCOV expirable.ChanOrValue[*proto.LocationInfo],
) {
	// try to get it as a station, returning a value if it is one
	// stations are locally-present static data
	info := r.TryAddStationById(locationId)
	if info != nil {
		infoCOV = expirable.NewCOVValue(expirable.New(info, fetch.MAX_EXPIRES))
	} else {
		chn := expirable.NewChanResult[*proto.LocationInfo](x.Ctx(), 1, 0)
		go expirable.P3Transceive(
			chn,
			x, r, locationId,
			protoGetStructureInfo,
		)
		infoCOV = expirable.NewCOVChan(chn)
	}
	return infoCOV
}

// get it as a structure
// structures are dynamic data, and fetched from ESI
func protoGetStructureInfo(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	locationId int64,
) (
	info *proto.LocationInfo,
	expires time.Time,
	err error,
) {
	structureInfo, expires, err := GetStructureInfo(x, locationId)
	if err != nil {
		info = nil
	} else if structureInfo == nil {
		info = r.AddUndefinedStructure(locationId)
	} else {
		info = r.AddStructure(
			locationId,
			structureInfo.Name,
			structureInfo.Forbidden,
			structureInfo.SolarSystemId,
		)
	}
	return info, expires, err
}
