package contractqueue

import (
	"time"

	"github.com/WiggidyW/etco-go/appraisal"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/contracts"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
	"github.com/WiggidyW/etco-go/protoregistry"
)

func ProtoGetBuybackContractQueue(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
) (
	entries []*proto.BuybackContractQueueEntry,
	expires time.Time,
	err error,
) {
	return protoGetContractQueue(
		x,
		r,
		contracts.GetBuybackContracts,
		appraisal.ProtoGetBuybackAppraisal,
		func(
			code string,
			contract *proto.Contract,
			appraisal *proto.BuybackAppraisal,
		) *proto.BuybackContractQueueEntry {
			return &proto.BuybackContractQueueEntry{
				Code:      code,
				Contract:  contract,
				Appraisal: appraisal,
			}
		},
	)
}

func ProtoGetShopContractQueue(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
) (
	entries []*proto.ShopContractQueueEntry,
	expires time.Time,
	err error,
) {
	return protoGetContractQueue(
		x,
		r,
		contracts.GetShopContracts,
		appraisal.ProtoGetShopAppraisal,
		func(
			code string,
			contract *proto.Contract,
			appraisal *proto.ShopAppraisal,
		) *proto.ShopContractQueueEntry {
			return &proto.ShopContractQueueEntry{
				Code:      code,
				Contract:  contract,
				Appraisal: appraisal,
			}
		},
	)
}

func ProtoGetHaulContractQueue(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
) (
	entries []*proto.HaulContractQueueEntry,
	expires time.Time,
	err error,
) {
	return protoGetContractQueue(
		x,
		r,
		contracts.GetHaulContracts,
		appraisal.ProtoGetHaulAppraisal,
		func(
			code string,
			contract *proto.Contract,
			appraisal *proto.HaulAppraisal,
		) *proto.HaulContractQueueEntry {
			return &proto.HaulContractQueueEntry{
				Code:      code,
				Contract:  contract,
				Appraisal: appraisal,
			}
		},
	)
}

type getContractsFunc func(
	x cache.Context,
) (
	contracts map[string]contracts.Contract,
	expires time.Time,
	err error,
)

type getAppraisalFunc[A any] func(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
	include_items bool,
) (
	appraisal A,
	expires time.Time,
	err error,
)

type newEntryFunc[A any, E any] func(
	code string,
	contract *proto.Contract,
	appraisal A,
) *E

func protoGetContractQueue[A proto.Nullable, E any](
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	getContracts getContractsFunc,
	getAppraisal getAppraisalFunc[A],
	newEntry newEntryFunc[A, E],
) (
	entries []*E,
	expires time.Time,
	err error,
) {
	// fetch contracts
	var contracts map[string]contracts.Contract
	contracts, expires, err = getContracts(x)
	if err != nil {
		return nil, expires, err
	}

	// fetch entries for each contract in a goroutine
	x, cancel := x.WithCancel()
	defer cancel()
	chn := expirable.NewChanResult[*E](x.Ctx(), len(contracts), 0)
	for code, contract := range contracts {
		go expirable.P6Transceive(
			chn,
			x, r, code, contract, getAppraisal, newEntry,
			protoGetContractQueueEntry,
		)
	}

	// recv entries
	entries = make([]*E, 0, len(contracts))
	var entry *E
	for i := 0; i < len(contracts); i++ {
		entry, expires, err = chn.RecvExpMin(expires)
		if err != nil {
			return nil, expires, err
		} else if entry != nil {
			entries = append(entries, entry)
		}
	}
	return entries, expires, nil
}

func protoGetContractQueueEntry[A proto.Nullable, E any](
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
	contract contracts.Contract,
	getAppraisal getAppraisalFunc[A],
	newEntry newEntryFunc[A, E],
) (
	entry *E,
	expires time.Time,
	err error,
) {
	x, cancel := x.WithCancel()
	defer cancel()

	// get the chanOrValue for location info (val if station, chan if structure)
	startLocationInfoCOVPtr, locationInfoCOV :=
		contracts.ProtoGetLocationInfoCOV(x, r, contract)

	// fetch the appraisal, returning nil if not found
	var appraisal A
	appraisal, expires, err = getAppraisal(x, r, code, false)
	if err != nil {
		if protoerr.ErrToProtoErr(err).Code != protoerr.NOT_FOUND {
			return nil, expires, err
		} else {
			return nil, expires, nil
		}
	} else if appraisal.IsNil() {
		return nil, expires, nil
	}

	// recv the location info
	var locationInfo *proto.LocationInfo
	locationInfo, expires, err = locationInfoCOV.RecvExpMin(expires)
	if err != nil {
		return nil, expires, err
	}

	// recv the start location info
	var startLocationInfo *proto.LocationInfo
	if startLocationInfoCOVPtr != nil {
		startLocationInfo, expires, err =
			startLocationInfoCOVPtr.RecvExpMin(expires)
		if err != nil {
			return nil, expires, err
		}
	}

	entry = newEntry(
		code,
		contract.ToProto(startLocationInfo, locationInfo),
		appraisal,
	)
	return entry, expires, nil
}
