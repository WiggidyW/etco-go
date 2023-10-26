package proto

import (
	"context"

	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go/client/contracts"
	"github.com/WiggidyW/etco-go/client/structureinfo"
	"github.com/WiggidyW/etco-go/kind"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBContractQueueParams struct {
	LocationInfoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker]
	StoreKind           kind.StoreKind
}

type LocationInfoChanAndResult struct {
	Channel  chanresult.ChanResult[*proto.LocationInfo]
	Result   *proto.LocationInfo
	Received bool
}

func NewLocationInfoChanAndResult(ctx context.Context) *LocationInfoChanAndResult {
	return &LocationInfoChanAndResult{
		Channel: chanresult.NewChanResult[*proto.LocationInfo](
			ctx,
			1,
			0,
		),
		// Result: nil,
		// Received: false,
	}
}

func (licr *LocationInfoChanAndResult) Recv() (
	info *proto.LocationInfo,
	err error,
) {
	if licr.Received {
		return licr.Result, nil
	} else {
		licr.Result, err = licr.Channel.Recv()
		if err != nil {
			return nil, err
		}
		licr.Received = true
		return licr.Result, nil
	}
}

type PBContractQueueClient struct {
	rContractsClient    contracts.WC_ContractsClient
	structureInfoClient structureinfo.WC_StructureInfoClient
}

func NewPBContractQueueClient(
	rContractsClient contracts.WC_ContractsClient,
	structureInfoClient structureinfo.WC_StructureInfoClient,
) PBContractQueueClient {
	return PBContractQueueClient{
		rContractsClient:    rContractsClient,
		structureInfoClient: structureInfoClient,
	}
}

func (cqc PBContractQueueClient) Fetch(
	ctx context.Context,
	params PBContractQueueParams,
) (
	entries []*proto.ContractQueueEntry,
	err error,
) {
	// fetch the current contracts for the requested store kind
	rContracts, err := cqc.rContractsClient.Fetch(
		ctx,
		contracts.ContractsParams{},
	)
	if err != nil {
		return nil, err
	}
	var rContractsMap map[string]contracts.Contract
	if params.StoreKind == kind.Shop {
		rContractsMap = rContracts.Data().ShopContracts
	} else {
		rContractsMap = rContracts.Data().BuybackContracts
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// map each location ID to a LocationInfoChanAndResult
	chnsLocationInfoMap := make(map[int64]*LocationInfoChanAndResult)
	for _, rContract := range rContractsMap {
		_, ok := chnsLocationInfoMap[rContract.LocationId]
		if ok {
			continue
		}
		// create a new one, and start a goroutine to fetch the location info
		chnAndResult := NewLocationInfoChanAndResult(ctx)
		chnsLocationInfoMap[rContract.LocationId] = chnAndResult
		go cqc.transceiveFetchLocationInfo(
			ctx,
			params.LocationInfoSession,
			rContract.LocationId,
			chnAndResult.Channel,
		)
	}

	// convert contracts to queue entries
	entries = make(
		[]*proto.ContractQueueEntry,
		0,
		len(rContractsMap),
	)
	for appraisalCode, rContract := range rContractsMap {
		// receive the location info
		chnAndResult := chnsLocationInfoMap[rContract.LocationId]
		locationInfo, err := chnAndResult.Recv()
		if err != nil {
			return nil, err
		}
		// append a new entry
		entries = append(entries, &proto.ContractQueueEntry{
			Code:                 appraisalCode,
			Contract:             protoutil.NewPBContract(rContract),
			ContractLocationInfo: locationInfo,
		})
	}

	return entries, nil
}

func (cqc PBContractQueueClient) transceiveFetchLocationInfo(
	ctx context.Context,
	infoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker],
	locationId int64,
	chnSend chanresult.ChanResult[*proto.LocationInfo],
) (err error) {
	locationInfo, err := cqc.fetchLocationInfo(
		ctx,
		infoSession,
		locationId,
	)
	if err != nil {
		return chnSend.SendErr(err)
	}
	return chnSend.SendOk(locationInfo)
}

func (cqc PBContractQueueClient) fetchLocationInfo(
	ctx context.Context,
	infoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker],
	locationId int64,
) (locationInfo *proto.LocationInfo, err error) {
	locationInfo, shouldFetchStructureInfo :=
		protoutil.MaybeGetExistingInfoOrTryAddAsStation(
			infoSession,
			locationId,
		)
	if shouldFetchStructureInfo {
		rStructureInfo, err := cqc.structureInfoClient.Fetch(
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
