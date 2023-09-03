package proto

import (
	"context"

	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go/client/appraisal"
	"github.com/WiggidyW/etco-go/client/contracts"
	"github.com/WiggidyW/etco-go/client/structureinfo"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
	"github.com/WiggidyW/etco-go/util"
)

// TODO: MAKE THIS LESS COMPLICATED, COMPLEX, AND CONFUSING
// crazy channel logic
// fairly fast but at what cost

type shopCodeAndNewAppraisals struct {
	codeAppraisal AppraisalWithCharacter[proto.ShopAppraisal]
	newAppraisal  *proto.ShopAppraisal
}

func (scana shopCodeAndNewAppraisals) Unwrap() (
	newAppraisal *proto.ShopAppraisal,
	codeAppraisal *proto.ShopAppraisal,
	characterId int32,
) {
	codeAppraisal, characterId = scana.codeAppraisal.Unwrap()
	return scana.newAppraisal, codeAppraisal, characterId
}

type codeLocationInfoChannels struct {
	chnLocationId   chanresult.ChanResult[int64]
	chnLocationInfo chanresult.ChanResult[*proto.LocationInfo]
}

func newCodeLocationInfoChannels(ctx context.Context) codeLocationInfoChannels {
	return codeLocationInfoChannels{
		chnLocationId: chanresult.
			NewChanResult[int64](ctx, 1, 0),
		chnLocationInfo: chanresult.
			NewChanResult[*proto.LocationInfo](ctx, 1, 0),
	}
}

type PBShopContractQueueClient struct {
	pbGetShopAppraisalClient PBGetShopAppraisalClient[*staticdb.SyncIndexMap]
	pbNewShopAppraisalClient PBNewShopAppraisalClient[*staticdb.SyncIndexMap]
	pbContractItemsClient    PBContractItemsClient[*staticdb.SyncIndexMap]
	rContractsClient         contracts.WC_ContractsClient
	structureInfoClient      structureinfo.WC_StructureInfoClient
}

func (gscqc PBShopContractQueueClient) Fetch(
	ctx context.Context,
	params PBContractQueueParams,
) (
	entries []*proto.ShopContractQueueEntry,
	err error,
) {
	rContracts, err := gscqc.rContractsClient.Fetch(
		ctx,
		contracts.ContractsParams{},
	)
	if err != nil {
		return entries, err
	}
	rShopContracts := rContracts.Data().ShopContracts

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnSendShopContractQueueEntry, chnRecvShopContractQueueEntry :=
		chanresult.NewChanResult[*proto.ShopContractQueueEntry](
			ctx,
			len(rShopContracts),
			0,
		).Split()
	chnsLocationInfoMap :=
		make(map[int64]*[]chanresult.ChanResult[*proto.LocationInfo])
	chnsCodeLocationInfos := make(
		[]codeLocationInfoChannels,
		0,
		len(rShopContracts),
	)

	// TODO: make this a function
	// For each contract, append a channel to the locationid->[]channel map
	// and start a goroutine to fetch the queue entry
	for appraisalCode, rContract := range rShopContracts {
		chnContractLocationInfo :=
			chanresult.NewChanResult[*proto.LocationInfo](ctx, 1, 0)
		chnsLocationInfo, ok :=
			chnsLocationInfoMap[rContract.LocationId]
		if ok {
			*chnsLocationInfo = append(
				*chnsLocationInfo,
				chnContractLocationInfo,
			)
		} else {
			chnsLocationInfo =
				&[]chanresult.ChanResult[*proto.LocationInfo]{
					chnContractLocationInfo,
				}
			chnsLocationInfoMap[rContract.LocationId] =
				chnsLocationInfo
		}

		chnsCodeLocationInfo := newCodeLocationInfoChannels(ctx)
		chnsCodeLocationInfos = append(
			chnsCodeLocationInfos,
			chnsCodeLocationInfo,
		)

		go gscqc.transceiveFetchEntry(
			ctx,
			params,
			appraisalCode,
			rContract,
			chnsCodeLocationInfo.chnLocationId.ToSend(),
			chnContractLocationInfo.ToRecv(),
			chnsCodeLocationInfo.chnLocationInfo.ToRecv(),
			chnSendShopContractQueueEntry,
		)
	}

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
		go gscqc.multiTransceiveFetchLocationInfo(
			ctx,
			params.LocationInfoSession,
			locationId,
			*chnsContractLocationInfo...,
		)
	}

	// finally, collect the queue entries
	entries = make(
		[]*proto.ShopContractQueueEntry,
		0,
		len(rShopContracts),
	)
	for i := 0; i < len(rShopContracts); i++ {
		entry, err := chnRecvShopContractQueueEntry.Recv()
		if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (gscqc PBShopContractQueueClient) multiTransceiveFetchLocationInfo(
	ctx context.Context,
	infoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker],
	locationId int64,
	variadicChnSendPB ...chanresult.ChanResult[*proto.LocationInfo],
) (err error) {
	var locationInfo *proto.LocationInfo
	if locationId == 0 {
		locationInfo = nil
	} else {
		locationInfo, err = gscqc.fetchLocationInfo(
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

func (gscqc PBShopContractQueueClient) fetchLocationInfo(
	ctx context.Context,
	infoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker],
	locationId int64,
) (locationInfo *proto.LocationInfo, err error) {
	locationInfo, shouldFetchStructureInfo := protoutil.
		MaybeGetExistingInfoOrTryAddAsStation(
			infoSession,
			locationId,
		)
	if shouldFetchStructureInfo {
		rStructureInfo, err := gscqc.structureInfoClient.Fetch(
			ctx,
			structureinfo.StructureInfoParams{
				StructureId: locationId,
			},
		)
		if err != nil {
			return nil, err
		}
		locationInfo = protoutil.MaybeAddStructureInfo(
			infoSession,
			locationId,
			rStructureInfo.Data().Forbidden,
			rStructureInfo.Data().Name,
			rStructureInfo.Data().SystemId,
		)
	}
	return locationInfo, nil
}

func (gscqc PBShopContractQueueClient) transceiveFetchEntry(
	ctx context.Context,
	params PBContractQueueParams,
	appraisalCode string,
	rContract contracts.Contract,
	chnSendLocationId chanresult.ChanSendResult[int64],
	chnRecvContractLocationInfo chanresult.ChanRecvResult[*proto.LocationInfo],
	chnRecvCodeLocationInfo chanresult.ChanRecvResult[*proto.LocationInfo],
	chnSend chanresult.ChanSendResult[*proto.ShopContractQueueEntry],
) error {
	pbShopContractQueueEntry, err := gscqc.fetchEntry(
		ctx,
		params,
		appraisalCode,
		rContract,
		chnSendLocationId,
		chnRecvContractLocationInfo,
		chnRecvCodeLocationInfo,
	)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(pbShopContractQueueEntry)
	}
}

func (gscqc PBShopContractQueueClient) fetchEntry(
	ctx context.Context,
	params PBContractQueueParams,
	appraisalCode string,
	rContract contracts.Contract,
	chnSendLocationId chanresult.ChanSendResult[int64],
	chnRecvContractLocationInfo chanresult.ChanRecvResult[*proto.LocationInfo],
	chnRecvCodeLocationInfo chanresult.ChanRecvResult[*proto.LocationInfo],
) (entry *proto.ShopContractQueueEntry, err error) {
	entry = &proto.ShopContractQueueEntry{}

	if params.QueueInclude == CQI_NONE {
		if err = chnSendLocationId.SendOk(0); err != nil {
			return nil, err
		}
	}

	if params.QueueInclude == CQI_ITEMS {
		if err = chnSendLocationId.SendOk(0); err != nil {
			return nil, err
		}
		if entry.ContractItems, err = gscqc.fetchContractItems(
			ctx,
			params.TypeNamingSession,
			rContract.ContractId,
		); err != nil {
			return nil, err
		}

	} else if params.QueueInclude == CQI_CODE_APPRAISAL {
		if entry.CodeAppraisal, entry.AppraisalCharacterId, err = util.
			Unwrap2WithErr(gscqc.fetchCodeAppraisal(
				ctx,
				params.TypeNamingSession,
				appraisalCode,
				chnSendLocationId,
			)); err != nil {
			return nil, err
		}

	} else if params.QueueInclude == CQI_ITEMS_AND_CODE_APPRAISAL {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		chnSendContractItems, chnRecvContractItems := chanresult.
			NewChanResult[[]*proto.ContractItem](ctx, 1, 0).Split()
		go gscqc.transceiveFetchContractItems(
			ctx,
			params.TypeNamingSession,
			rContract.ContractId,
			chnSendContractItems,
		)

		if entry.CodeAppraisal,
			entry.AppraisalCharacterId,
			err = util.Unwrap2WithErr(
			gscqc.fetchCodeAppraisal(
				ctx,
				params.TypeNamingSession,
				appraisalCode,
				chnSendLocationId,
			),
		); err != nil {
			return nil, err
		}

		entry.ContractItems, err = chnRecvContractItems.Recv()
		if err != nil {
			return nil, err
		}

	} else if params.QueueInclude == CQI_CODE_APPRAISAL_AND_NEW_APPRAISAL {
		if entry.NewAppraisal,
			entry.CodeAppraisal,
			entry.AppraisalCharacterId,
			err = util.Unwrap3WithErr(
			gscqc.fetchCodeAndNewAppraisals(
				ctx,
				params.TypeNamingSession,
				appraisalCode,
				chnSendLocationId,
			),
		); err != nil {
			return nil, err
		}

	} else { // if params.ShopQueueInclude == BQI_ITEMS_AND_CODE_APPRAISAL_AND_NEW_APPRAISAL {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		chnSendContractItems, chnRecvContractItems := chanresult.
			NewChanResult[[]*proto.ContractItem](ctx, 1, 0).Split()
		go gscqc.transceiveFetchContractItems(
			ctx,
			params.TypeNamingSession,
			rContract.ContractId,
			chnSendContractItems,
		)

		if entry.NewAppraisal,
			entry.CodeAppraisal,
			entry.AppraisalCharacterId,
			err = util.Unwrap3WithErr(
			gscqc.fetchCodeAndNewAppraisals(
				ctx,
				params.TypeNamingSession,
				appraisalCode,
				chnSendLocationId,
			),
		); err != nil {
			return nil, err
		}

		entry.ContractItems, err = chnRecvContractItems.Recv()
		if err != nil {
			return nil, err
		}
	}

	entry.Contract = protoutil.NewPBContract(rContract)

	entry.ContractLocationInfo, err = chnRecvContractLocationInfo.Recv()
	if err != nil {
		return nil, err
	}

	entry.AppraisalLocationInfo, err = chnRecvCodeLocationInfo.Recv()
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (gscqc PBShopContractQueueClient) transceiveFetchContractItems(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	contractId int32,
	chnSendPB chanresult.ChanSendResult[[]*proto.ContractItem],
) error {
	pbContractItems, err := gscqc.fetchContractItems(
		ctx,
		namingSesssion,
		contractId,
	)
	if err != nil {
		return chnSendPB.SendErr(err)
	} else {
		return chnSendPB.SendOk(pbContractItems)
	}
}

func (gscqc PBShopContractQueueClient) fetchContractItems(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	contractId int32,
) (
	pbContractItems []*proto.ContractItem,
	err error,
) {
	return gscqc.pbContractItemsClient.Fetch(
		ctx,
		PBContractItemsParams[*staticdb.SyncIndexMap]{
			TypeNamingSession: namingSesssion,
			ContractId:        contractId,
		},
	)
}

// func (gscqc PBShopContractQueueClient[
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
// 	appraisals, err := gscqc.fetchCodeAndNewAppraisals(
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

func (gscqc PBShopContractQueueClient) fetchCodeAndNewAppraisals(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	appraisalCode string,
	// required channel to send the locationID for code location info
	chnSendLocationId chanresult.ChanSendResult[int64],
) (
	appraisals shopCodeAndNewAppraisals,
	err error,
) {
	appraisals.codeAppraisal, err = gscqc.fetchCodeAppraisal(
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

	appraisals.newAppraisal, err = gscqc.fetchNewAppraisal(
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

// func (gscqc PBShopContractQueueClient) transceiveFetchCodeAppraisal(
// 	ctx context.Context,
// 	namingSesssion *staticdb.TypeNamingSession[IM],
// 	appraisalCode string,
// 	// required channel to send the locationID for code location info
// 	chnSendLocationId chanresult.ChanSendResult[int64],
// 	chnSend chanresult.ChanSendResult[AppraisalWithCharacter[proto.ShopAppraisal]],
// ) error {
// 	shopAppraisal, err := gscqc.fetchCodeAppraisal(
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

func (gscqc PBShopContractQueueClient) fetchCodeAppraisal(
	ctx context.Context,
	namingSession *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	appraisalCode string,
	// required channel to send the locationID for code location info
	chnSendLocationId chanresult.ChanSendResult[int64],
) (
	appraisal AppraisalWithCharacter[proto.ShopAppraisal],
	err error,
) {
	appraisal, err = gscqc.pbGetShopAppraisalClient.Fetch(
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

// func (gscqc PBShopContractQueueClient) transceiveFetchNewAppraisal(
// 	ctx context.Context,
// 	namingSesssion *staticdb.TypeNamingSession[IM],
// 	rItems []appraisal.BasicItem,
// 	locationId int64,
// 	chnSend chanresult.ChanSendResult[*proto.ShopAppraisal],
// ) error {
// 	pbShopAppraisal, err := gscqc.fetchNewAppraisal(
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

func (gscqc PBShopContractQueueClient) fetchNewAppraisal(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	rItems []appraisal.BasicItem,
	locationId int64,
) (*proto.ShopAppraisal, error) {
	return gscqc.pbNewShopAppraisalClient.Fetch(
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
