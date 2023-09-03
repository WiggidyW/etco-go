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

type buybackCodeAndNewAppraisals struct {
	codeAppraisal AppraisalWithCharacter[proto.BuybackAppraisal]
	newAppraisal  *proto.BuybackAppraisal
}

func (bcana buybackCodeAndNewAppraisals) Unwrap() (
	newAppraisal *proto.BuybackAppraisal,
	codeAppraisal *proto.BuybackAppraisal,
	characterId int32,
) {
	codeAppraisal, characterId = bcana.codeAppraisal.Unwrap()
	return bcana.newAppraisal, codeAppraisal, characterId
}

type PBBuybackContractQueueClient struct {
	pbGetBuybackAppraisalClient PBGetBuybackAppraisalClient[*staticdb.SyncIndexMap]
	pbNewBuybackAppraisalClient PBNewBuybackAppraisalClient[*staticdb.SyncIndexMap]
	pbContractItemsClient       PBContractItemsClient[*staticdb.SyncIndexMap]
	rContractsClient            contracts.WC_ContractsClient
	structureInfoClient         structureinfo.WC_StructureInfoClient
}

func (gbcqc PBBuybackContractQueueClient) Fetch(
	ctx context.Context,
	params PBContractQueueParams,
) (
	entries []*proto.BuybackContractQueueEntry,
	err error,
) {
	rContracts, err := gbcqc.rContractsClient.Fetch(ctx, contracts.ContractsParams{})
	if err != nil {
		return entries, err
	}
	rBuybackContracts := rContracts.Data().BuybackContracts

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnSendBuybackContractQueueEntry, chnRecvBuybackContractQueueEntry :=
		chanresult.NewChanResult[*proto.BuybackContractQueueEntry](
			ctx,
			len(rBuybackContracts),
			0,
		).Split()
	chnsLocationInfoMap :=
		make(map[int64]*[]chanresult.ChanResult[*proto.LocationInfo])

	// TODO: make this a function
	// For each contract, append a channel to the locationid->[]channel map
	// and start a goroutine to fetch the queue entry
	for appraisalCode, rContract := range rBuybackContracts {
		chnLocationInfo := chanresult.NewChanResult[*proto.LocationInfo](ctx, 1, 0)
		if chnsLocationInfo, ok := chnsLocationInfoMap[rContract.LocationId]; ok {
			*chnsLocationInfo = append(*chnsLocationInfo, chnLocationInfo)
		} else {
			chnsLocationInfo = &[]chanresult.ChanResult[*proto.LocationInfo]{
				chnLocationInfo,
			}
			chnsLocationInfoMap[rContract.LocationId] = chnsLocationInfo
		}
		go gbcqc.transceiveFetchEntry(
			ctx,
			params,
			appraisalCode,
			rContract,
			chnLocationInfo.ToRecv(),
			chnSendBuybackContractQueueEntry,
		)
	}

	// For each location, start a goroutine to fetch the location info
	// and send it to all the channels in the locationid->[]channel map
	for locationId, chnsLocationInfo := range chnsLocationInfoMap {
		// TODO: We can actually do this without waiting!
		go gbcqc.multiTransceiveFetchLocationInfo(
			ctx,
			params.LocationInfoSession,
			locationId,
			*chnsLocationInfo...,
		)
	}

	// finally, collect the queue entries
	entries = make(
		[]*proto.BuybackContractQueueEntry,
		0,
		len(rBuybackContracts),
	)
	for i := 0; i < len(rBuybackContracts); i++ {
		entry, err := chnRecvBuybackContractQueueEntry.Recv()
		if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (gbcqc PBBuybackContractQueueClient) multiTransceiveFetchLocationInfo(
	ctx context.Context,
	infoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker],
	locationId int64,
	variadicChnSendPB ...chanresult.ChanResult[*proto.LocationInfo],
) (err error) {
	locationInfo, err := gbcqc.fetchLocationInfo(
		ctx,
		infoSession,
		locationId,
	)
	if err != nil {
		return err
	}
	for _, chnSendPB := range variadicChnSendPB {
		err = chnSendPB.SendOk(locationInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gbcqc PBBuybackContractQueueClient) fetchLocationInfo(
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
		rStructureInfo, err := gbcqc.structureInfoClient.Fetch(
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

func (gbcqc PBBuybackContractQueueClient) transceiveFetchEntry(
	ctx context.Context,
	params PBContractQueueParams,
	appraisalCode string,
	rContract contracts.Contract,
	chnRecvLocationInfo chanresult.ChanRecvResult[*proto.LocationInfo],
	chnSend chanresult.ChanSendResult[*proto.BuybackContractQueueEntry],
) error {
	pbBuybackContractQueueEntry, err := gbcqc.fetchEntry(
		ctx,
		params,
		appraisalCode,
		rContract,
		chnRecvLocationInfo,
	)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(pbBuybackContractQueueEntry)
	}
}

func (gbcqc PBBuybackContractQueueClient) fetchEntry(
	ctx context.Context,
	params PBContractQueueParams,
	appraisalCode string,
	rContract contracts.Contract,
	chnRecvLocationInfo chanresult.ChanRecvResult[*proto.LocationInfo],
) (entry *proto.BuybackContractQueueEntry, err error) {
	entry = &proto.BuybackContractQueueEntry{}

	if params.QueueInclude == CQI_ITEMS {
		if entry.ContractItems, err = gbcqc.fetchContractItems(
			ctx,
			params.TypeNamingSession,
			rContract.ContractId,
		); err != nil {
			return nil, err
		}

	} else if params.QueueInclude == CQI_CODE_APPRAISAL {
		entry.CodeAppraisal, entry.AppraisalCharacterId, err = util.
			Unwrap2WithErr(gbcqc.fetchCodeAppraisal(
				ctx,
				params.TypeNamingSession,
				appraisalCode,
			))
		if err != nil {
			return nil, err
		}

	} else if params.QueueInclude == CQI_ITEMS_AND_CODE_APPRAISAL {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		chnSendContractItems, chnRecvContractItems := chanresult.
			NewChanResult[[]*proto.ContractItem](ctx, 1, 0).Split()
		go gbcqc.transceiveFetchContractItems(
			ctx,
			params.TypeNamingSession,
			rContract.ContractId,
			chnSendContractItems,
		)

		if entry.CodeAppraisal,
			entry.AppraisalCharacterId,
			err = util.Unwrap2WithErr(
			gbcqc.fetchCodeAppraisal(
				ctx,
				params.TypeNamingSession,
				appraisalCode,
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
			gbcqc.fetchCodeAndNewAppraisals(
				ctx,
				params.TypeNamingSession,
				appraisalCode,
			),
		); err != nil {
			return nil, err
		}

	} else { // if params.BuybackQueueInclude == BQI_ITEMS_AND_CODE_APPRAISAL_AND_NEW_APPRAISAL {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		chnSendContractItems, chnRecvContractItems := chanresult.
			NewChanResult[[]*proto.ContractItem](ctx, 1, 0).Split()
		go gbcqc.transceiveFetchContractItems(
			ctx,
			params.TypeNamingSession,
			rContract.ContractId,
			chnSendContractItems,
		)

		if entry.NewAppraisal,
			entry.CodeAppraisal,
			entry.AppraisalCharacterId,
			err = util.Unwrap3WithErr(
			gbcqc.fetchCodeAndNewAppraisals(
				ctx,
				params.TypeNamingSession,
				appraisalCode,
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

	entry.ContractLocationInfo, err = chnRecvLocationInfo.Recv()
	if err != nil {
		return nil, err
	}

	// set the appraisal region id only if include_code_appraisal + include_location_info
	if entry.CodeAppraisal != nil {
		entry.AppraisalRegionId = protoutil.MaybeAddSystem(
			params.LocationInfoSession,
			entry.CodeAppraisal.SystemId,
		)
	}

	return entry, nil
}

func (gbcqc PBBuybackContractQueueClient) transceiveFetchContractItems(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	contractId int32,
	chnSend chanresult.ChanSendResult[[]*proto.ContractItem],
) error {
	pbContractItems, err := gbcqc.fetchContractItems(
		ctx,
		namingSesssion,
		contractId,
	)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(pbContractItems)
	}
}

func (gbcqc PBBuybackContractQueueClient) fetchContractItems(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	contractId int32,
) (
	pbContractItems []*proto.ContractItem,
	err error,
) {
	return gbcqc.pbContractItemsClient.Fetch(
		ctx,
		PBContractItemsParams[*staticdb.SyncIndexMap]{
			TypeNamingSession: namingSesssion,
			ContractId:        contractId,
		},
	)
}

// func (gbcqc PBBuybackContractQueueClient[
// 	IM,
// 	LN,
// ]) transceiveFetchCodeAndNewAppraisals(
// 	ctx context.Context,
// 	namingSesssion *staticdb.TypeNamingSession[IM],
// 	appraisalCode string,
// 	chnSend chanresult.ChanSendResult[buybackCodeAndNewAppraisals],
// ) error {
// 	appraisals, err := gbcqc.fetchCodeAndNewAppraisals(
// 		ctx,
// 		namingSesssion,
// 		appraisalCode,
// 	)
// 	if err != nil {
// 		return chnSend.SendErr(err)
// 	} else {
// 		return chnSend.SendOk(appraisals)
// 	}
// }

func (gbcqc PBBuybackContractQueueClient) fetchCodeAndNewAppraisals(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	appraisalCode string,
) (
	appraisals buybackCodeAndNewAppraisals,
	err error,
) {
	appraisals.codeAppraisal, err = gbcqc.fetchCodeAppraisal(
		ctx,
		namingSesssion,
		appraisalCode,
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

	appraisals.newAppraisal, err = gbcqc.fetchNewAppraisal(
		ctx,
		namingSesssion,
		rItems,
		appraisals.codeAppraisal.Appraisal.SystemId,
	)
	if err != nil {
		return appraisals, err
	}

	return appraisals, nil
}

// func (gbcqc PBBuybackContractQueueClient) transceiveFetchCodeAppraisal(
// 	ctx context.Context,
// 	namingSesssion *staticdb.TypeNamingSession[IM],
// 	appraisalCode string,
// 	chnSend chanresult.ChanSendResult[AppraisalWithCharacter[proto.BuybackAppraisal]],
// ) error {
// 	buybackAppraisal, err := gbcqc.fetchCodeAppraisal(
// 		ctx,
// 		namingSesssion,
// 		appraisalCode,
// 	)
// 	if err != nil {
// 		return chnSend.SendErr(err)
// 	} else {
// 		return chnSend.SendOk(buybackAppraisal)
// 	}
// }

func (gbcqc PBBuybackContractQueueClient) fetchCodeAppraisal(
	ctx context.Context,
	namingSession *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	appraisalCode string,
) (
	appraisal AppraisalWithCharacter[proto.BuybackAppraisal],
	err error,
) {
	return gbcqc.pbGetBuybackAppraisalClient.Fetch(
		ctx,
		PBGetAppraisalParams[*staticdb.SyncIndexMap]{
			TypeNamingSession: namingSession,
			AppraisalCode:     appraisalCode,
		},
	)
}

// func (gbcqc PBBuybackContractQueueClient) transceiveFetchNewAppraisal(
// 	ctx context.Context,
// 	namingSesssion *staticdb.TypeNamingSession[IM],
// 	rItems []appraisal.BasicItem,
// 	systemId int32,
// 	chnSend chanresult.ChanSendResult[*proto.BuybackAppraisal],
// ) error {
// 	appraisal, err := gbcqc.fetchNewAppraisal(
// 		ctx,
// 		namingSesssion,
// 		rItems,
// 		systemId,
// 	)
// 	if err != nil {
// 		return chnSend.SendErr(err)
// 	} else {
// 		return chnSend.SendOk(appraisal)
// 	}
// }

func (gbcqc PBBuybackContractQueueClient) fetchNewAppraisal(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	rItems []appraisal.BasicItem,
	systemId int32,
) (
	*proto.BuybackAppraisal,
	error,
) {
	return gbcqc.pbNewBuybackAppraisalClient.Fetch(
		ctx,
		PBNewBuybackAppraisalParams[*staticdb.SyncIndexMap]{
			TypeNamingSession: namingSesssion,
			Items:             rItems,
			SystemId:          systemId,
			CharacterId:       nil,
			Save:              false,
		},
	)
}
