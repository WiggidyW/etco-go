package proto

import (
	"context"

	"github.com/WiggidyW/chanresult"

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

func (sbac PBStatusShopAppraisalClient) Fetch(
	ctx context.Context,
	params PBStatusAppraisalParams,
) (
	partialStatus PartialStatusShopAppraisal,
	err error,
) {
	rShopQueueRep, err := sbac.rShopQueueClient.Fetch(
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
				partialStatus.InPurchaseQueue = true
				return partialStatus, nil
			}
		}
	}

	// if params.StatusInclude == ASI_NONE {}

	if params.StatusInclude == ASI_ITEMS {
		partialStatus.ContractItems, err = sbac.fetchContractItems(
			ctx,
			params.TypeNamingSession,
			rContract.ContractId,
		)
		if err != nil {
			return partialStatus, err
		}

	} else if params.StatusInclude == ASI_LOCATION_INFO {
		partialStatus.LocationInfo, err = sbac.fetchLocationInfo(
			ctx,
			params.LocationInfoSession,
			rContract.LocationId,
		)
		if err != nil {
			return partialStatus, err
		}

		// } else if params.StatusInclude == ASI_ITEMS_AND_LOCATION_INFO {
	} else {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		chnSendContractItems, chnRecvContractItems := chanresult.
			NewChanResult[[]*proto.ContractItem](ctx, 1, 0).Split()
		go sbac.transceiveFetchContractItems(
			ctx,
			params.TypeNamingSession,
			rContract.ContractId,
			chnSendContractItems,
		)

		partialStatus.LocationInfo, err = sbac.fetchLocationInfo(
			ctx,
			params.LocationInfoSession,
			rContract.LocationId,
		)
		if err != nil {
			return partialStatus, err
		}

		partialStatus.ContractItems, err = chnRecvContractItems.Recv()
		if err != nil {
			return partialStatus, err
		}
	}

	partialStatus.Contract = protoutil.NewPBContract(rContract)

	return partialStatus, nil
}

func (sbac PBStatusShopAppraisalClient) fetchLocationInfo(
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

func (sbac PBStatusShopAppraisalClient) transceiveFetchContractItems(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.LocalIndexMap],
	contractId int32,
	chnSend chanresult.ChanSendResult[[]*proto.ContractItem],
) error {
	pbContractItems, err := sbac.fetchContractItems(
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

func (sbac PBStatusShopAppraisalClient) fetchContractItems(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.LocalIndexMap],
	contractId int32,
) (
	pbContractItems []*proto.ContractItem,
	err error,
) {
	return sbac.pbContractItemsClient.Fetch(
		ctx,
		PBContractItemsParams[*staticdb.LocalIndexMap]{
			TypeNamingSession: namingSesssion,
			ContractId:        contractId,
		},
	)
}
