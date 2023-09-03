package proto

import (
	"context"

	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go/client/appraisal"
	"github.com/WiggidyW/etco-go/client/shopqueue"
	"github.com/WiggidyW/etco-go/client/structureinfo"
	"github.com/WiggidyW/etco-go/proto"
	pu "github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
	"github.com/WiggidyW/etco-go/util"
)

// TODO: MAKE THIS LESS COMPLICATED, COMPLEX, AND CONFUSING
// crazy channel logic
// fairly fast but at what cost

type PBShopPurchaseQueueClient struct {
	pbGetShopAppraisalClient PBGetShopAppraisalClient[*staticdb.SyncIndexMap]
	pbNewShopAppraisalClient PBNewShopAppraisalClient[*staticdb.SyncIndexMap]
	rShopQueueClient         shopqueue.ShopQueueClient
	structureInfoClient      structureinfo.WC_StructureInfoClient
}

func (gspqc PBShopPurchaseQueueClient) Fetch(
	ctx context.Context,
	params PBPurchaseQueueParams,
) (
	entries []*proto.ShopPurchaseQueueEntry,
	err error,
) {
	rShopQueueRep, err := gspqc.rShopQueueClient.Fetch(
		ctx,
		shopqueue.ShopQueueParams{},
	)
	if err != nil {
		return entries, err
	}
	rShopQueue := rShopQueueRep.ParsedShopQueue

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnSendShopPurchaseQueueEntry, chnRecvShopPurchaseQueueEntry :=
		chanresult.NewChanResult[*proto.ShopPurchaseQueueEntry](
			ctx,
			len(rShopQueue),
			0,
		).Split()
	chnsCodeLocationInfos := make(
		[]codeLocationInfoChannels,
		0,
		len(rShopQueue),
	)

	// TODO: make this a function
	// For each contract, append a channel to the locationid->[]channel map
	// and start a goroutine to fetch the queue entry
	for _, appraisalCode := range rShopQueue {
		chnsCodeLocationInfo := newCodeLocationInfoChannels(ctx)
		chnsCodeLocationInfos = append(
			chnsCodeLocationInfos,
			chnsCodeLocationInfo,
		)

		go gspqc.transceiveFetchEntry(
			ctx,
			params,
			appraisalCode,
			chnsCodeLocationInfo.chnLocationId.ToSend(),
			chnsCodeLocationInfo.chnLocationInfo.ToRecv(),
			chnSendShopPurchaseQueueEntry,
		)
	}

	chnsLocationInfoMap :=
		make(map[int64]*[]chanresult.ChanResult[*proto.LocationInfo])

	// for each entry thread, receive the code appraisals location ID and
	// append the codeLocationInfo channel to the locationid->[]channel map
	for _, chnsCodeLocationInfo := range chnsCodeLocationInfos {
		locationId, err := chnsCodeLocationInfo.chnLocationId.Recv()
		if err != nil {
			return entries, err
		}
		chnsLocationInfo, ok := chnsLocationInfoMap[locationId]
		if ok {
			*chnsLocationInfo = append(
				*chnsLocationInfo,
				chnsCodeLocationInfo.chnLocationInfo,
			)
		} else {
			chnsLocationInfo =
				&[]chanresult.ChanResult[*proto.LocationInfo]{
					chnsCodeLocationInfo.chnLocationInfo,
				}
			chnsLocationInfoMap[locationId] =
				chnsLocationInfo
		}
	}

	// For each location, start a goroutine to fetch the location info
	// and send it to all the channels in the locationid->[]channel map
	for locationId, chnsContractLocationInfo := range chnsLocationInfoMap {
		go gspqc.multiTransceiveFetchLocationInfo(
			ctx,
			params.LocationInfoSession,
			locationId,
			*chnsContractLocationInfo...,
		)
	}

	// finally, collect the queue entries
	entries = make(
		[]*proto.ShopPurchaseQueueEntry,
		0,
		len(rShopQueue),
	)
	for i := 0; i < len(rShopQueue); i++ {
		entry, err := chnRecvShopPurchaseQueueEntry.Recv()
		if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (gspqc PBShopPurchaseQueueClient) multiTransceiveFetchLocationInfo(
	ctx context.Context,
	infoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker],
	locationId int64,
	variadicChnSendPB ...chanresult.ChanResult[*proto.LocationInfo],
) (err error) {
	var locationInfo *proto.LocationInfo
	if locationId == 0 {
		locationInfo = nil
	} else {
		locationInfo, err = gspqc.fetchLocationInfo(
			ctx,
			infoSession,
			locationId,
		)
		if err != nil {
			return err
		}
	}

	for _, chnSendPB := range variadicChnSendPB {
		err = chnSendPB.SendOk(locationInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gspqc PBShopPurchaseQueueClient) fetchLocationInfo(
	ctx context.Context,
	infoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker],
	locationId int64,
) (locationInfo *proto.LocationInfo, err error) {
	locationInfo, shouldFetchStructureInfo := pu.
		MaybeGetExistingInfoOrTryAddAsStation(
			infoSession,
			locationId,
		)
	if shouldFetchStructureInfo {
		rStructureInfo, err := gspqc.structureInfoClient.Fetch(
			ctx,
			structureinfo.StructureInfoParams{
				StructureId: locationId,
			},
		)
		if err != nil {
			return nil, err
		}
		locationInfo = pu.MaybeAddStructureInfo(
			infoSession,
			locationId,
			rStructureInfo.Data().Forbidden,
			rStructureInfo.Data().Name,
			rStructureInfo.Data().SystemId,
		)
	}
	return locationInfo, nil
}

func (gspqc PBShopPurchaseQueueClient) transceiveFetchEntry(
	ctx context.Context,
	params PBPurchaseQueueParams,
	appraisalCode string,
	chnSendLocationId chanresult.ChanSendResult[int64],
	chnRecvCodeLocationInfo chanresult.ChanRecvResult[*proto.LocationInfo],
	chnSend chanresult.ChanSendResult[*proto.ShopPurchaseQueueEntry],
) error {
	pbShopPurchaseQueueEntry, err := gspqc.fetchEntry(
		ctx,
		params,
		appraisalCode,
		chnSendLocationId,
		chnRecvCodeLocationInfo,
	)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(pbShopPurchaseQueueEntry)
	}
}

func (gspqc PBShopPurchaseQueueClient) fetchEntry(
	ctx context.Context,
	params PBPurchaseQueueParams,
	appraisalCode string,
	chnSendLocationId chanresult.ChanSendResult[int64],
	chnRecvCodeLocationInfo chanresult.ChanRecvResult[*proto.LocationInfo],
) (entry *proto.ShopPurchaseQueueEntry, err error) {
	entry = &proto.ShopPurchaseQueueEntry{}

	if params.QueueInclude == PQI_NONE {
		if err = chnSendLocationId.SendOk(0); err != nil {
			return nil, err
		}

	} else if params.QueueInclude == PQI_CODE_APPRAISAL {
		if entry.CodeAppraisal, entry.AppraisalCharacterId, err = util.
			Unwrap2WithErr(gspqc.fetchCodeAppraisal(
				ctx,
				params.TypeNamingSession,
				appraisalCode,
				chnSendLocationId,
			)); err != nil {
			return nil, err
		}

	} else { // if params.QueueInclude == PQI_CODE_APPRAISAL_AND_NEW_APPRAISAL
		if entry.NewAppraisal,
			entry.CodeAppraisal,
			entry.AppraisalCharacterId,
			err = util.Unwrap3WithErr(
			gspqc.fetchCodeAndNewAppraisals(
				ctx,
				params.TypeNamingSession,
				appraisalCode,
				chnSendLocationId,
			),
		); err != nil {
			return nil, err
		}

	}

	entry.AppraisalLocationInfo, err = chnRecvCodeLocationInfo.Recv()
	if err != nil {
		return nil, err
	}

	return entry, nil
}

// func (gspqc PBShopPurchaseQueueClient[
// 	IM,
// 	LN,
// ]) transceiveFetchCodeAndNewAppraisals(
// 	ctx context.Context,
// 	namingSesssion *staticdb.TypeNamingSession[IM],
// 	appraisalCode string,
// 	// required channel to send the locationID for code location info
// 	chnSendLocationId chanresult.ChanSendResult[int64],
// 	chnSend chanresult.ChanSendResult[shopCodeAndNewAppraisals],
// ) error {
// 	appraisals, err := gspqc.fetchCodeAndNewAppraisals(
// 		ctx,
// 		namingSesssion,
// 		appraisalCode,
// 		chnSendLocationId,
// 	)
// 	if err != nil {
// 		return chnSend.SendErr(err)
// 	} else {
// 		return chnSend.SendOk(appraisals)
// 	}
// }

func (gspqc PBShopPurchaseQueueClient) fetchCodeAndNewAppraisals(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	appraisalCode string,
	// required channel to send the locationID for code location info
	chnSendLocationId chanresult.ChanSendResult[int64],
) (
	appraisals shopCodeAndNewAppraisals,
	err error,
) {
	appraisals.codeAppraisal, err = gspqc.fetchCodeAppraisal(
		ctx,
		namingSesssion,
		appraisalCode,
		chnSendLocationId,
	)
	if err != nil {
		return appraisals, err
	}

	rItems := make(
		[]appraisal.BasicItem,
		0,
		len(appraisals.codeAppraisal.Appraisal.Items),
	)
	for _, item := range appraisals.codeAppraisal.Appraisal.Items {
		rItems = append(rItems, appraisal.BasicItem{
			TypeId:   item.TypeId,
			Quantity: item.Quantity,
		})
	}

	appraisals.newAppraisal, err = gspqc.fetchNewAppraisal(
		ctx,
		namingSesssion,
		rItems,
		appraisals.codeAppraisal.Appraisal.LocationId,
	)
	if err != nil {
		return appraisals, err
	}

	return appraisals, nil
}

// func (gspqc PBShopPurchaseQueueClient) transceiveFetchCodeAppraisal(
// 	ctx context.Context,
// 	namingSesssion *staticdb.TypeNamingSession[IM],
// 	appraisalCode string,
// 	// required channel to send the locationID for code location info
// 	chnSendLocationId chanresult.ChanSendResult[int64],
// 	chnSend chanresult.ChanSendResult[AppraisalWithCharacter[proto.ShopAppraisal]],
// ) error {
// 	shopAppraisal, err := gspqc.fetchCodeAppraisal(
// 		ctx,
// 		namingSesssion,
// 		appraisalCode,
// 		chnSendLocationId,
// 	)
// 	if err != nil {
// 		return chnSend.SendErr(err)
// 	} else {
// 		return chnSend.SendOk(shopAppraisal)
// 	}
// }

func (gspqc PBShopPurchaseQueueClient) fetchCodeAppraisal(
	ctx context.Context,
	namingSession *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	appraisalCode string,
	// required channel to send the locationID for code location info
	chnSendLocationId chanresult.ChanSendResult[int64],
) (
	appraisal AppraisalWithCharacter[proto.ShopAppraisal],
	err error,
) {
	appraisal, err = gspqc.pbGetShopAppraisalClient.Fetch(
		ctx,
		PBGetAppraisalParams[*staticdb.SyncIndexMap]{
			TypeNamingSession: namingSession,
			AppraisalCode:     appraisalCode,
		},
	)
	if err != nil {
		return appraisal, err
	}

	err = chnSendLocationId.SendOk(appraisal.Appraisal.LocationId)
	if err != nil {
		return appraisal, err
	}

	return appraisal, nil
}

// func (gspqc PBShopPurchaseQueueClient) transceiveFetchNewAppraisal(
// 	ctx context.Context,
// 	namingSesssion *staticdb.TypeNamingSession[IM],
// 	rItems []appraisal.BasicItem,
// 	locationId int64,
// 	chnSend chanresult.ChanSendResult[*proto.ShopAppraisal],
// ) error {
// 	pbShopAppraisal, err := gspqc.fetchNewAppraisal(
// 		ctx,
// 		namingSesssion,
// 		rItems,
// 		locationId,
// 	)
// 	if err != nil {
// 		return chnSend.SendErr(err)
// 	} else {
// 		return chnSend.SendOk(pbShopAppraisal)
// 	}
// }

func (gspqc PBShopPurchaseQueueClient) fetchNewAppraisal(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	rItems []appraisal.BasicItem,
	locationId int64,
) (*proto.ShopAppraisal, error) {
	return gspqc.pbNewShopAppraisalClient.Fetch(
		ctx,
		PBNewShopAppraisalParams[*staticdb.SyncIndexMap]{
			TypeNamingSession: namingSesssion,
			Items:             rItems,
			LocationId:        locationId,
			CharacterId:       0,
			IncludeCode:       false,
		},
	)
}
