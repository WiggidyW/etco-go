package proto

import (
	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type ShopLocationsParams struct {
	LocationInfoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker]
}

type PBShopLocationsClient struct{}

func NewPBShopLocationsClient() PBShopLocationsClient {
	return PBShopLocationsClient{}
}

func (slc PBShopLocationsClient) Fetch(
	x cache.Context,
	params ShopLocationsParams,
) (
	rep []*proto.ShopLocation,
	err error,
) {
	UNSAFE_ShopLocations := staticdb.UnsafeGetCoreShopLocations()
	rep = make([]*proto.ShopLocation, 0, len(UNSAFE_ShopLocations))

	// if we don't need location info, convert it to PB and return now
	if params.LocationInfoSession == nil {
		for locationId := range UNSAFE_ShopLocations {
			rep = append(rep, &proto.ShopLocation{
				LocationId: locationId,
				// LocationInfo: nil,
			})
		}
		return rep, nil
	}

	x, cancel := x.WithCancel()
	defer cancel()

	// send out a location info fetch for each location ID
	chnSendLocationInfo, chnRecvLocationInfo := chanresult.
		NewChanResult[LocationInfoWithLocationId](x.Ctx(), 1, 0).Split()
	for locationId := range UNSAFE_ShopLocations {
		go slc.transceiveFetchLocationInfo(
			x,
			locationId,
			params.LocationInfoSession,
			chnSendLocationInfo,
		)
	}

	// receive all location info and append
	for i := 0; i < len(UNSAFE_ShopLocations); i++ {
		locationInfoWithId, err := chnRecvLocationInfo.Recv()
		if err != nil {
			return nil, err
		}
		rep = append(rep, &proto.ShopLocation{
			LocationId:   locationInfoWithId.LocationId,
			LocationInfo: locationInfoWithId.LocationInfo,
		})
	}

	return rep, nil
}

func (slc PBShopLocationsClient) transceiveFetchLocationInfo(
	x cache.Context,
	locationid int64,
	infoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker],
	chnSend chanresult.ChanSendResult[LocationInfoWithLocationId],
) error {
	locationInfoWithId, err := slc.fetchLocationInfo(
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

func (slc PBShopLocationsClient) fetchLocationInfo(
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
