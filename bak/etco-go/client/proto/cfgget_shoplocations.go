package proto

import (
	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type LocationInfoWithLocationId struct {
	LocationId   int64
	LocationInfo *proto.LocationInfo
}

type PartialCfgShopLocationsResponse struct {
	Locations       map[int64]*proto.CfgShopLocation
	LocationInfoMap map[int64]*proto.LocationInfo
}

type CfgGetShopLocationsParams struct {
	LocationInfoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker]
}

type CfgGetShopLocationsClient struct{}

func NewCfgGetShopLocationsClient() CfgGetShopLocationsClient {
	return CfgGetShopLocationsClient{}
}

func (gslc CfgGetShopLocationsClient) Fetch(
	x cache.Context,
	params CfgGetShopLocationsParams,
) (
	rep *PartialCfgShopLocationsResponse,
	err error,
) {
	// fetch web shop locations
	webShopLocations, err := gslc.fetchWebShopLocations(x)
	if err != nil {
		return nil, err
	}

	// if we don't need location info, convert it to PB and return now
	if params.LocationInfoSession == nil {
		return &PartialCfgShopLocationsResponse{
			Locations: protoutil.NewPBCfgShopLocations(
				webShopLocations,
			),
			// LocationInfoMap: nil,
		}, nil
	}

	x, cancel := x.WithCancel()
	defer cancel()

	// send out a location info fetch for each location ID
	chnSendLocationInfo, chnRecvLocationInfo := chanresult.
		NewChanResult[LocationInfoWithLocationId](x.Ctx(), 1, 0).Split()
	for locationId := range webShopLocations {
		go gslc.transceiveFetchLocationInfo(
			x,
			locationId,
			params.LocationInfoSession,
			chnSendLocationInfo,
		)
	}

	// initialize response
	rep = &PartialCfgShopLocationsResponse{
		Locations: protoutil.NewPBCfgShopLocations(webShopLocations),
		LocationInfoMap: make(
			map[int64]*proto.LocationInfo,
			len(webShopLocations),
		),
	}

	// receive all location info and insert to location info map
	for i := 0; i < len(webShopLocations); i++ {
		locationInfoWithId, err := chnRecvLocationInfo.Recv()
		if err != nil {
			return nil, err
		}
		rep.LocationInfoMap[locationInfoWithId.LocationId] =
			locationInfoWithId.LocationInfo
	}

	return rep, nil
}

func (gslc CfgGetShopLocationsClient) transceiveFetchLocationInfo(
	x cache.Context,
	locationid int64,
	infoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker],
	chnSend chanresult.ChanSendResult[LocationInfoWithLocationId],
) error {
	locationInfoWithId, err := gslc.fetchLocationInfo(
		x,
		locationid,
		infoSession,
	)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(locationInfoWithId)
	}
}

func (gslc CfgGetShopLocationsClient) fetchLocationInfo(
	x cache.Context,
	locationId int64,
	infoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker],
) (
	locationInfoWithId LocationInfoWithLocationId,
	err error,
) {
	locationInfo, shouldFetchStructureInfo := protoutil.
		MaybeGetExistingInfoOrTryAddAsStation(infoSession, locationId)

	if !shouldFetchStructureInfo {
		return LocationInfoWithLocationId{
			LocationId:   locationId,
			LocationInfo: locationInfo,
		}, nil
	}

	structureInfo, _, err := esi.GetStructureInfo( // TODO: Handle Nil (it never happens atm)
		x,
		locationId,
	)
	if err != nil {
		return locationInfoWithId, err
	}
	return LocationInfoWithLocationId{
		LocationId: locationId,
		LocationInfo: protoutil.MaybeAddStructureInfo(
			infoSession,
			locationId,
			structureInfo.Forbidden,
			structureInfo.Name,
			structureInfo.SolarSystemId,
		),
	}, nil
}

func (gslc CfgGetShopLocationsClient) fetchWebShopLocations(
	x cache.Context,
) (
	shopLocations map[b.LocationId]b.WebShopLocation,
	err error,
) {
	shopLocations, _, err = bucket.GetWebShopLocations(x)
	return shopLocations, err
}
