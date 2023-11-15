package proto

import (
	"context"

	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/contracts"
	"github.com/WiggidyW/etco-go/esi"
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

type PBContractQueueClient struct{}

func NewPBContractQueueClient() PBContractQueueClient {
	return PBContractQueueClient{}
}

func (cqc PBContractQueueClient) Fetch(
	x cache.Context,
	params PBContractQueueParams,
) (
	entries []*proto.ContractQueueEntry,
	err error,
) {
	// fetch the current contracts for the requested store kind
	rContracts, _, err := contracts.GetContracts(x)
	if err != nil {
		return nil, err
	}
	var rContractsMap map[string]contracts.Contract
	if params.StoreKind == kind.Shop {
		rContractsMap = rContracts.ShopContracts
	} else {
		rContractsMap = rContracts.BuybackContracts
	}

	x, cancel := x.WithCancel()
	defer cancel()

	// map each location ID to a LocationInfoChanAndResult
	chnsLocationInfoMap := make(map[int64]*LocationInfoChanAndResult)
	for _, rContract := range rContractsMap {
		_, ok := chnsLocationInfoMap[rContract.LocationId]
		if ok {
			continue
		}
		// create a new one, and start a goroutine to fetch the location info
		chnAndResult := NewLocationInfoChanAndResult(x.Ctx())
		chnsLocationInfoMap[rContract.LocationId] = chnAndResult
		go cqc.transceiveFetchLocationInfo(
			x,
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
	x cache.Context,
	infoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker],
	locationId int64,
	chnSend chanresult.ChanResult[*proto.LocationInfo],
) (err error) {
	locationInfo, err := cqc.fetchLocationInfo(
		x,
		infoSession,
		locationId,
	)
	if err != nil {
		return chnSend.SendErr(err)
	}
	return chnSend.SendOk(locationInfo)
}

func (cqc PBContractQueueClient) fetchLocationInfo(
	x cache.Context,
	infoSession *staticdb.LocationInfoSession[*staticdb.SyncLocationNamerTracker],
	locationId int64,
) (locationInfo *proto.LocationInfo, err error) {
	locationInfo, shouldFetchStructureInfo :=
		protoutil.MaybeGetExistingInfoOrTryAddAsStation(
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
