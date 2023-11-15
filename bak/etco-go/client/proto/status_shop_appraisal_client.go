package proto

import (
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/contracts"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/purchasequeue"
	"github.com/WiggidyW/etco-go/staticdb"

	"github.com/WiggidyW/chanresult"
)

type PartialStatusShopAppraisal struct {
	Contract        *proto.Contract
	InPurchaseQueue bool
	ContractItems   []*proto.ContractItem
	LocationInfo    *proto.LocationInfo
}

type PBStatusShopAppraisalClient struct {
	pbContractItemsClient PBContractItemsClient[*staticdb.LocalIndexMap]
}

func NewPBStatusShopAppraisalClient(
	pbContractItemsClient PBContractItemsClient[*staticdb.LocalIndexMap],
) PBStatusShopAppraisalClient {
	return PBStatusShopAppraisalClient{pbContractItemsClient}
}

func (ssac PBStatusShopAppraisalClient) Fetch(
	x cache.Context,
	params PBStatusAppraisalParams,
) (
	partialStatus PartialStatusShopAppraisal,
	err error,
) {
	queueX, queueCancel := x.WithCancel()
	chnQueue := expirable.NewChanResult[map[int64][]string](queueX.Ctx(), 1, 0)
	go expirable.Param1Transceive(
		chnQueue,
		x,
		purchasequeue.GetPurchaseQueue,
	)

	rShopContracts, _, err := contracts.GetShopContracts(x)
	if err != nil {
		return partialStatus, err
	}

	// get the contract we're requested to fetch
	rContract, ok := rShopContracts[params.AppraisalCode]
	if !ok {
		defer queueCancel()
		queue, _, err := chnQueue.RecvExp()
		if err != nil {
			return partialStatus, err
		}
		for _, codes := range queue {
			for _, code := range codes {
				if code == params.AppraisalCode {
					// no contract found + in purchase queue
					partialStatus.InPurchaseQueue = true
					return partialStatus, nil
				}
			}
		}
		// no contract found + not in purchase queue
		return partialStatus, nil
	} else {
		queueCancel()
	}

	x, cancel := x.WithCancel()
	defer cancel()

	// send a goroutine to fetch the contract items
	chnSendContractItems, chnRecvContractItems := chanresult.
		NewChanResult[[]*proto.ContractItem](x.Ctx(), 1, 0).Split()
	go ssac.transceiveFetchContractItems(
		x,
		params.TypeNamingSession,
		rContract.ContractId,
		params.IncludeItems && rContract.Status != contracts.Deleted,
		chnSendContractItems,
	)

	// fetch location info
	partialStatus.LocationInfo, err = ssac.fetchLocationInfo(
		x,
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
	x cache.Context,
	infoSession *staticdb.LocationInfoSession[*staticdb.LocalLocationNamerTracker],
	locationId int64,
) (locationInfo *proto.LocationInfo, err error) {
	locationInfo, shouldFetchStructureInfo := protoutil.
		MaybeGetExistingInfoOrTryAddAsStation[*staticdb.LocalLocationNamerTracker](
		infoSession,
		locationId,
	)
	if shouldFetchStructureInfo {
		rStructureInfo, _, err := esi.GetStructureInfo( // TODO: Handle Nil (it never happens atm)
			x,
			locationId,
		)
		if err != nil {
			return nil, err
		}
		locationInfo = protoutil.MaybeAddStructureInfo(
			infoSession,
			locationId,
			rStructureInfo.Forbidden,
			rStructureInfo.Name,
			rStructureInfo.SolarSystemId,
		)
	}
	return locationInfo, nil
}

func (ssac PBStatusShopAppraisalClient) transceiveFetchContractItems(
	x cache.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.LocalIndexMap],
	contractId int32,
	includeItems bool,
	chnSend chanresult.ChanSendResult[[]*proto.ContractItem],
) error {
	pbContractItems, err := ssac.fetchContractItems(
		x,
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
	x cache.Context,
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
			x,
			PBContractItemsParams[*staticdb.LocalIndexMap]{
				TypeNamingSession: namingSesssion,
				ContractId:        contractId,
			},
		)
	}
}
