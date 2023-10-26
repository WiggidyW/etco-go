package proto

import (
	"context"

	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go/client/contracts"
	"github.com/WiggidyW/etco-go/client/structureinfo"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PartialStatusBuybackAppraisal struct {
	Contract      *proto.Contract
	ContractItems []*proto.ContractItem
	LocationInfo  *proto.LocationInfo
}

type PBStatusBuybackAppraisalClient struct {
	pbContractItemsClient PBContractItemsClient[*staticdb.LocalIndexMap]
	rContractsClient      contracts.WC_ContractsClient
	structureInfoClient   structureinfo.WC_StructureInfoClient
}

func NewPBStatusBuybackAppraisalClient(
	pbContractItemsClient PBContractItemsClient[*staticdb.LocalIndexMap],
	rContractsClient contracts.WC_ContractsClient,
	structureInfoClient structureinfo.WC_StructureInfoClient,
) PBStatusBuybackAppraisalClient {
	return PBStatusBuybackAppraisalClient{
		pbContractItemsClient,
		rContractsClient,
		structureInfoClient,
	}
}

func (sbac PBStatusBuybackAppraisalClient) Fetch(
	ctx context.Context,
	params PBStatusAppraisalParams,
) (
	partialStatus PartialStatusBuybackAppraisal,
	err error,
) {
	rContracts, err := sbac.rContractsClient.Fetch(
		ctx,
		contracts.ContractsParams{},
	)
	if err != nil {
		return partialStatus, err
	}

	// get the contract we're requested to fetch
	rContract, ok := rContracts.Data().
		BuybackContracts[params.AppraisalCode]
	if !ok {
		// no contract found
		return partialStatus, nil
	}

	// send a goroutine to fetch the contract items
	chnSendContractItems, chnRecvContractItems := chanresult.
		NewChanResult[[]*proto.ContractItem](ctx, 1, 0).Split()
	go sbac.transceiveFetchContractItems(
		ctx,
		params.TypeNamingSession,
		rContract.ContractId,
		params.IncludeItems && rContract.Status != contracts.Deleted,
		chnSendContractItems,
	)

	// fetch location info
	partialStatus.LocationInfo, err = sbac.fetchLocationInfo(
		ctx,
		params.LocationInfoSession,
		rContract.LocationId,
	)
	if err != nil {
		return partialStatus, err
	}

	// set contract
	partialStatus.Contract = protoutil.NewPBContract(rContract)

	// receive contract items
	partialStatus.ContractItems, err = chnRecvContractItems.Recv()
	if err != nil {
		return partialStatus, err
	}

	return partialStatus, nil
}

func (sbac PBStatusBuybackAppraisalClient) fetchLocationInfo(
	ctx context.Context,
	infoSession *staticdb.LocationInfoSession[*staticdb.LocalLocationNamerTracker],
	locationId int64,
) (locationInfo *proto.LocationInfo, err error) {
	locationInfo, shouldFetchStructureInfo :=
		protoutil.MaybeGetExistingInfoOrTryAddAsStation(
			infoSession,
			locationId,
		)
	if shouldFetchStructureInfo {
		rStructureInfo, err := sbac.structureInfoClient.Fetch(
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

func (sbac PBStatusBuybackAppraisalClient) transceiveFetchContractItems(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.LocalIndexMap],
	contractId int32,
	includeItems bool,
	chnSend chanresult.ChanSendResult[[]*proto.ContractItem],
) error {
	pbContractItems, err := sbac.fetchContractItems(
		ctx,
		namingSesssion,
		contractId,
		includeItems,
	)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(pbContractItems)
	}
}

func (sbac PBStatusBuybackAppraisalClient) fetchContractItems(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.LocalIndexMap],
	contractId int32,
	includeItems bool,
) (
	pbContractItems []*proto.ContractItem,
	err error,
) {
	if !includeItems {
		return nil, nil
	} else {
		return sbac.pbContractItemsClient.Fetch(
			ctx,
			PBContractItemsParams[*staticdb.LocalIndexMap]{
				TypeNamingSession: namingSesssion,
				ContractId:        contractId,
			},
		)
	}
}
