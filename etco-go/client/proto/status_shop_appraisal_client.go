package proto

import (
	"context"

	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go/client/contracts"
	"github.com/WiggidyW/etco-go/client/shopqueue"
	"github.com/WiggidyW/etco-go/client/structureinfo"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PartialStatusShopAppraisal struct {
	Contract        *proto.Contract
	InPurchaseQueue bool
	ContractItems   []*proto.ContractItem
	LocationInfo    *proto.LocationInfo
}

type PBStatusShopAppraisalClient struct {
	pbContractItemsClient PBContractItemsClient[*staticdb.LocalIndexMap]
	rShopQueueClient      shopqueue.ShopQueueClient
	structureInfoClient   structureinfo.WC_StructureInfoClient
}

func NewPBStatusShopAppraisalClient(
	pbContractItemsClient PBContractItemsClient[*staticdb.LocalIndexMap],
	rShopQueueClient shopqueue.ShopQueueClient,
	structureInfoClient structureinfo.WC_StructureInfoClient,
) PBStatusShopAppraisalClient {
	return PBStatusShopAppraisalClient{
		pbContractItemsClient,
		rShopQueueClient,
		structureInfoClient,
	}
}

func (ssac PBStatusShopAppraisalClient) Fetch(
	ctx context.Context,
	params PBStatusAppraisalParams,
) (
	partialStatus PartialStatusShopAppraisal,
	err error,
) {
	rShopQueueRep, err := ssac.rShopQueueClient.Fetch(
		ctx,
		shopqueue.ShopQueueParams{},
	)
	if err != nil {
		return partialStatus, err
	}

	// get the contract we're requested to fetch
	rContract, ok := rShopQueueRep.ShopContracts[params.AppraisalCode]
	if !ok {
		for _, code := range rShopQueueRep.ParsedShopQueue {
			if code == params.AppraisalCode {
				// no contract found + in purchase queue
				partialStatus.InPurchaseQueue = true
				return partialStatus, nil
			}
		}
		// no contract found + not in purchase queue
		return partialStatus, nil
	}

	// send a goroutine to fetch the contract items
	chnSendContractItems, chnRecvContractItems := chanresult.
		NewChanResult[[]*proto.ContractItem](ctx, 1, 0).Split()
	go ssac.transceiveFetchContractItems(
		ctx,
		params.TypeNamingSession,
		rContract.ContractId,
		params.IncludeItems && rContract.Status != contracts.Deleted,
		chnSendContractItems,
	)

	// fetch location info
	partialStatus.LocationInfo, err = ssac.fetchLocationInfo(
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

func (ssac PBStatusShopAppraisalClient) fetchLocationInfo(
	ctx context.Context,
	infoSession *staticdb.LocationInfoSession[*staticdb.LocalLocationNamerTracker],
	locationId int64,
) (locationInfo *proto.LocationInfo, err error) {
	locationInfo, shouldFetchStructureInfo := protoutil.
		MaybeGetExistingInfoOrTryAddAsStation[*staticdb.LocalLocationNamerTracker](
		infoSession,
		locationId,
	)
	if shouldFetchStructureInfo {
		rStructureInfo, err := ssac.structureInfoClient.Fetch(
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

func (ssac PBStatusShopAppraisalClient) transceiveFetchContractItems(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.LocalIndexMap],
	contractId int32,
	includeItems bool,
	chnSend chanresult.ChanSendResult[[]*proto.ContractItem],
) error {
	pbContractItems, err := ssac.fetchContractItems(
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

func (ssac PBStatusShopAppraisalClient) fetchContractItems(
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
		return ssac.pbContractItemsClient.Fetch(
			ctx,
			PBContractItemsParams[*staticdb.LocalIndexMap]{
				TypeNamingSession: namingSesssion,
				ContractId:        contractId,
			},
		)
	}
}
