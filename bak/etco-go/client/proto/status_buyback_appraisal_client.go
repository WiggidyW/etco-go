package proto

import (
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/contracts"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"

	"github.com/WiggidyW/chanresult"
)

type PartialStatusBuybackAppraisal struct {
	Contract      *proto.Contract
	ContractItems []*proto.ContractItem
	LocationInfo  *proto.LocationInfo
}

type PBStatusBuybackAppraisalClient struct {
	pbContractItemsClient PBContractItemsClient[*staticdb.LocalIndexMap]
}

func NewPBStatusBuybackAppraisalClient(
	pbContractItemsClient PBContractItemsClient[*staticdb.LocalIndexMap],
) PBStatusBuybackAppraisalClient {
	return PBStatusBuybackAppraisalClient{
		pbContractItemsClient,
	}
}

func (sbac PBStatusBuybackAppraisalClient) Fetch(
	x cache.Context,
	params PBStatusAppraisalParams,
) (
	partialStatus PartialStatusBuybackAppraisal,
	err error,
) {
	rContracts, _, err := contracts.GetBuybackContracts(x)
	if err != nil {
		return partialStatus, err
	}

	// get the contract we're requested to fetch
	rContract, ok := rContracts[params.AppraisalCode]
	if !ok {
		// no contract found
		return partialStatus, nil
	}

	x, cancel := x.WithCancel()
	defer cancel()

	// send a goroutine to fetch the contract items
	chnSendContractItems, chnRecvContractItems := chanresult.
		NewChanResult[[]*proto.ContractItem](x.Ctx(), 1, 0).Split()
	go sbac.transceiveFetchContractItems(
		x,
		params.TypeNamingSession,
		rContract.ContractId,
		params.IncludeItems && rContract.Status != contracts.Deleted,
		chnSendContractItems,
	)

	// fetch location info
	partialStatus.LocationInfo, err = sbac.fetchLocationInfo(
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

func (sbac PBStatusBuybackAppraisalClient) fetchLocationInfo(
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

func (sbac PBStatusBuybackAppraisalClient) transceiveFetchContractItems(
	x cache.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.LocalIndexMap],
	contractId int32,
	includeItems bool,
	chnSend chanresult.ChanSendResult[[]*proto.ContractItem],
) error {
	pbContractItems, err := sbac.fetchContractItems(
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

func (sbac PBStatusBuybackAppraisalClient) fetchContractItems(
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
		return sbac.pbContractItemsClient.Fetch(
			x,
			PBContractItemsParams[*staticdb.LocalIndexMap]{
				TypeNamingSession: namingSesssion,
				ContractId:        contractId,
			},
		)
	}
}
