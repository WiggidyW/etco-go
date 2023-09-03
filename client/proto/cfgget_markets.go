package proto

import (
	"context"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/client/bucket"
	"github.com/WiggidyW/etco-go/client/structureinfo"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PartialCfgMarketsResponse struct {
	Markets         map[string]*proto.CfgMarket
	LocationInfoMap map[int64]*proto.LocationInfo
}

type CfgGetMarketsParams struct {
	LocationInfoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker]
}

type CfgGetMarketsClient struct {
	webMarketsReaderClient bucket.WebMarketsReaderClient
	structureInfoClient    structureinfo.WC_StructureInfoClient
}

func (gslc CfgGetMarketsClient) Fetch(
	ctx context.Context,
	params CfgGetMarketsParams,
) (
	rep *PartialCfgMarketsResponse,
	err error,
) {
	// fetch web shop locations
	webMarkets, err := gslc.fetchWebMarkets(ctx)
	if err != nil {
		return nil, err
	}

	// if we don't need location info, convert it to PB and return now
	if params.LocationInfoSession == nil {
		return &PartialCfgMarketsResponse{
			Markets: protoutil.NewPBCfgMarkets(
				webMarkets,
			),
			// LocationInfoMap: nil,
		}, nil
	}

	// track unique location IDs and send out a fetch for each one

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// send out a location info fetch for each unique location ID
	uniqueLocationIds := make(map[int64]struct{}, len(webMarkets))
	chnSendLocationInfo, chnRecvLocationInfo := chanresult.
		NewChanResult[LocationInfoWithLocationId](ctx, 1, 0).Split()
	for _, webMarket := range webMarkets {
		if _, ok := uniqueLocationIds[webMarket.LocationId]; ok {
			continue
		}
		go gslc.transceiveFetchLocationInfo(
			ctx,
			webMarket.LocationId,
			params.LocationInfoSession,
			chnSendLocationInfo,
		)
		uniqueLocationIds[webMarket.LocationId] = struct{}{}
	}

	// initialize response
	rep = &PartialCfgMarketsResponse{
		Markets: protoutil.NewPBCfgMarkets(webMarkets),
		LocationInfoMap: make(
			map[int64]*proto.LocationInfo,
			len(uniqueLocationIds),
		),
	}

	// receive all location info and insert to location info map
	for i := 0; i < len(uniqueLocationIds); i++ {
		locationInfoWithId, err := chnRecvLocationInfo.Recv()
		if err != nil {
			return nil, err
		}
		rep.LocationInfoMap[locationInfoWithId.LocationId] =
			locationInfoWithId.LocationInfo
	}

	return rep, nil
}

func (gslc CfgGetMarketsClient) transceiveFetchLocationInfo(
	ctx context.Context,
	locationid int64,
	infoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker],
	chnSend chanresult.ChanSendResult[LocationInfoWithLocationId],
) error {
	locationInfoWithId, err := gslc.fetchLocationInfo(
		ctx,
		locationid,
		infoSession,
	)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(locationInfoWithId)
	}
}

func (gslc CfgGetMarketsClient) fetchLocationInfo(
	ctx context.Context,
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

	structureInfo, err := gslc.structureInfoClient.Fetch(
		ctx,
		structureinfo.StructureInfoParams{StructureId: locationId},
	)
	if err != nil {
		return locationInfoWithId, err
	}

	return LocationInfoWithLocationId{
		LocationId: locationId,
		LocationInfo: protoutil.MaybeAddStructureInfo(
			infoSession,
			locationId,
			structureInfo.Data().Forbidden,
			structureInfo.Data().Name,
			structureInfo.Data().SystemId,
		),
	}, nil
}

func (gslc CfgGetMarketsClient) fetchWebMarkets(
	ctx context.Context,
) (
	webMarkets map[b.MarketName]b.WebMarket,
	err error,
) {
	webMarketsRep, err := gslc.webMarketsReaderClient.Fetch(
		ctx,
		bucket.WebMarketsReaderParams{},
	)
	if err != nil {
		return nil, err
	} else {
		return webMarketsRep.Data(), nil
	}
}
