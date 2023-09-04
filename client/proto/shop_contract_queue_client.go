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

// TODO: DEDUPLICATE / DRY
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

type PBShopContractQueueClient struct {
	pbGetShopAppraisalClient PBGetShopAppraisalClient[*staticdb.SyncIndexMap]
	pbNewShopAppraisalClient PBNewShopAppraisalClient[*staticdb.SyncIndexMap]
	pbContractItemsClient    PBContractItemsClient[*staticdb.SyncIndexMap]
	rContractsClient         contracts.WC_ContractsClient
	structureInfoClient      structureinfo.WC_StructureInfoClient
}

func NewPBShopContractQueueClient(
	pbGetShopAppraisalClient PBGetShopAppraisalClient[*staticdb.SyncIndexMap],
	pbNewShopAppraisalClient PBNewShopAppraisalClient[*staticdb.SyncIndexMap],
	pbContractItemsClient PBContractItemsClient[*staticdb.SyncIndexMap],
	rContractsClient contracts.WC_ContractsClient,
	structureInfoClient structureinfo.WC_StructureInfoClient,
) PBShopContractQueueClient {
	return PBShopContractQueueClient{
		pbGetShopAppraisalClient: pbGetShopAppraisalClient,
		pbNewShopAppraisalClient: pbNewShopAppraisalClient,
		pbContractItemsClient:    pbContractItemsClient,
		rContractsClient:         rContractsClient,
		structureInfoClient:      structureInfoClient,
	}
}

func (scqc PBShopContractQueueClient) Fetch(
	ctx context.Context,
	params PBContractQueueParams,
) (
	entries []*proto.ShopContractQueueEntry,
	err error,
) {
	rContracts, err := scqc.rContractsClient.Fetch(
		ctx,
		contracts.ContractsParams{},
	)
	if err != nil {
		return entries, err
	}
	rShopContracts := rContracts.Data().ShopContracts

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnSendQueueEntry, chnRecvQueueEntry :=
		chanresult.NewChanResult[*proto.ShopContractQueueEntry](
			ctx,
			len(rShopContracts),
			0,
		).Split()
	chnsLocationInfoMap :=
		make(map[int64]*[]chanresult.ChanResult[*proto.LocationInfo])

	// TODO: make this a function
	// For each contract, append a channel to the locationid->[]channel map
	// and start a goroutine to fetch the queue entry
	for appraisalCode, rContract := range rShopContracts {
		chnLocationInfo :=
			chanresult.NewChanResult[*proto.LocationInfo](ctx, 1, 0)
		chnsLocationInfo, ok :=
			chnsLocationInfoMap[rContract.LocationId]
		if ok {
			*chnsLocationInfo = append(
				*chnsLocationInfo,
				chnLocationInfo,
			)
		} else {
			chnsLocationInfo =
				&[]chanresult.ChanResult[*proto.LocationInfo]{
					chnLocationInfo,
				}
			chnsLocationInfoMap[rContract.LocationId] =
				chnsLocationInfo
		}

		go scqc.transceiveFetchEntry(
			ctx,
			params,
			appraisalCode,
			rContract,
			chnLocationInfo.ToRecv(),
			chnSendQueueEntry,
		)
	}

	// For each location, start a goroutine to fetch the location info
	// and send it to all the channels in the locationid->[]channel map
	for locationId, chnsContractLocationInfo := range chnsLocationInfoMap {
		go scqc.multiTransceiveFetchLocationInfo(
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
		entry, err := chnRecvQueueEntry.Recv()
		if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (scqc PBShopContractQueueClient) multiTransceiveFetchLocationInfo(
	ctx context.Context,
	infoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker],
	locationId int64,
	chnsSend ...chanresult.ChanResult[*proto.LocationInfo],
) (err error) {
	locationInfo, err := scqc.fetchLocationInfo(
		ctx,
		infoSession,
		locationId,
	)
	if err != nil {
		return err
	}
	for _, chnSend := range chnsSend {
		err = chnSend.SendOk(locationInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (scqc PBShopContractQueueClient) fetchLocationInfo(
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
		rStructureInfo, err := scqc.structureInfoClient.Fetch(
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

func (scqc PBShopContractQueueClient) transceiveFetchEntry(
	ctx context.Context,
	params PBContractQueueParams,
	appraisalCode string,
	rContract contracts.Contract,
	chnRecvLocationInfo chanresult.ChanRecvResult[*proto.LocationInfo],
	chnSend chanresult.ChanSendResult[*proto.ShopContractQueueEntry],
) error {
	pbQueueEntry, err := scqc.fetchEntry(
		ctx,
		params,
		appraisalCode,
		rContract,
		chnRecvLocationInfo,
	)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(pbQueueEntry)
	}
}

func (scqc PBShopContractQueueClient) fetchEntry(
	ctx context.Context,
	params PBContractQueueParams,
	appraisalCode string,
	rContract contracts.Contract,
	chnRecvLocationInfo chanresult.ChanRecvResult[*proto.LocationInfo],
) (entry *proto.ShopContractQueueEntry, err error) {
	entry = &proto.ShopContractQueueEntry{Code: appraisalCode}

	// if params.QueueInclude == CQI_NONE {}

	if params.QueueInclude == CQI_ITEMS {
		if entry.ContractItems, err = scqc.fetchContractItems(
			ctx,
			params.TypeNamingSession,
			rContract.ContractId,
		); err != nil {
			return nil, err
		}

	} else if params.QueueInclude == CQI_CODE_APPRAISAL {
		if entry.CodeAppraisal, entry.AppraisalCharacterId, err = util.
			Unwrap2WithErr(scqc.fetchCodeAppraisal(
				ctx,
				params.TypeNamingSession,
				appraisalCode,
			)); err != nil {
			return nil, err
		}

	} else if params.QueueInclude == CQI_ITEMS_AND_CODE_APPRAISAL {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		chnSendContractItems, chnRecvContractItems := chanresult.
			NewChanResult[[]*proto.ContractItem](ctx, 1, 0).Split()
		go scqc.transceiveFetchContractItems(
			ctx,
			params.TypeNamingSession,
			rContract.ContractId,
			chnSendContractItems,
		)

		if entry.CodeAppraisal,
			entry.AppraisalCharacterId,
			err = util.Unwrap2WithErr(
			scqc.fetchCodeAppraisal(
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
			scqc.fetchCodeAndNewAppraisals(
				ctx,
				params.TypeNamingSession,
				appraisalCode,
			),
		); err != nil {
			return nil, err
		}

	} else if params.QueueInclude == CQI_ITEMS_AND_CODE_APPRAISAL_AND_NEW_APPRAISAL {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		chnSendContractItems, chnRecvContractItems := chanresult.
			NewChanResult[[]*proto.ContractItem](ctx, 1, 0).Split()
		go scqc.transceiveFetchContractItems(
			ctx,
			params.TypeNamingSession,
			rContract.ContractId,
			chnSendContractItems,
		)

		if entry.NewAppraisal,
			entry.CodeAppraisal,
			entry.AppraisalCharacterId,
			err = util.Unwrap3WithErr(
			scqc.fetchCodeAndNewAppraisals(
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

	return entry, nil
}

func (scqc PBShopContractQueueClient) transceiveFetchContractItems(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	contractId int32,
	chnSendPB chanresult.ChanSendResult[[]*proto.ContractItem],
) error {
	pbContractItems, err := scqc.fetchContractItems(
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

func (scqc PBShopContractQueueClient) fetchContractItems(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	contractId int32,
) (
	pbContractItems []*proto.ContractItem,
	err error,
) {
	return scqc.pbContractItemsClient.Fetch(
		ctx,
		PBContractItemsParams[*staticdb.SyncIndexMap]{
			TypeNamingSession: namingSesssion,
			ContractId:        contractId,
		},
	)
}

func (scqc PBShopContractQueueClient) fetchCodeAndNewAppraisals(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	appraisalCode string,
) (
	appraisals shopCodeAndNewAppraisals,
	err error,
) {
	appraisals.codeAppraisal, err = scqc.fetchCodeAppraisal(
		ctx,
		namingSesssion,
		appraisalCode,
	)
	if err != nil {
		return appraisals, err
	}

	appraisals.newAppraisal, err = scqc.fetchNewAppraisal(
		ctx,
		namingSesssion,
		protoutil.NewRBasicItems(
			appraisals.codeAppraisal.Appraisal.Items,
		),
		appraisals.codeAppraisal.Appraisal.LocationId,
	)
	if err != nil {
		return appraisals, err
	}

	return appraisals, nil
}

func (scqc PBShopContractQueueClient) fetchCodeAppraisal(
	ctx context.Context,
	namingSession *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	appraisalCode string,
) (
	appraisal AppraisalWithCharacter[proto.ShopAppraisal],
	err error,
) {
	return scqc.pbGetShopAppraisalClient.Fetch(
		ctx,
		PBGetAppraisalParams[*staticdb.SyncIndexMap]{
			TypeNamingSession: namingSession,
			AppraisalCode:     appraisalCode,
		},
	)
}

func (scqc PBShopContractQueueClient) fetchNewAppraisal(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	rItems []appraisal.BasicItem,
	locationId int64,
) (*proto.ShopAppraisal, error) {
	return scqc.pbNewShopAppraisalClient.Fetch(
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
