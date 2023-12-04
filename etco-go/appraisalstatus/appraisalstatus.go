package appraisalstatus

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/contracts"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"
	"github.com/WiggidyW/etco-go/purchasequeue"
)

type ProtoAppraisalStatusRep struct {
	Status        proto.AppraisalStatus
	Contract      *proto.Contract
	ContractItems []*proto.NamedBasicItem
}

func protoGetAppraisalStatus(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
	include_items bool,
	getStatusWithItems func(
		cache.Context,
		*protoregistry.ProtoRegistry,
		string,
	) (
		contracts.ProtoContractWithItemsRep,
		time.Time,
		error,
	),
	getStatus func(
		cache.Context,
		*protoregistry.ProtoRegistry,
		string,
	) (
		*proto.Contract,
		time.Time,
		error,
	),
) (
	rep ProtoAppraisalStatusRep,
	expires time.Time,
	err error,
) {
	var contractRep *ProtoAppraisalStatusRep
	if include_items {
		contractRep, expires, err = protoGetContractAppraisalStatusWithItems(
			x,
			r,
			code,
			getStatusWithItems,
		)
	} else {
		contractRep, expires, err = protoGetContractAppraisalStatus(
			x,
			r,
			code,
			getStatus,
		)
	}

	if err == nil && contractRep != nil {
		rep = *contractRep
	} else {
		rep = ProtoAppraisalStatusRep{
			Status:        proto.AppraisalStatus_AS_UNDEFINED,
			Contract:      nil,
			ContractItems: nil,
		}
	}

	return rep, expires, err
}

func ProtoGetBuybackAppraisalStatus(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
	include_items bool,
) (
	rep ProtoAppraisalStatusRep,
	expires time.Time,
	err error,
) {
	return protoGetAppraisalStatus(
		x,
		r,
		code,
		include_items,
		contracts.ProtoGetBuybackContractWithItems,
		contracts.ProtoGetBuybackContract,
	)
}

func ProtoGetHaulAppraisalStatus(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
	include_items bool,
) (
	rep ProtoAppraisalStatusRep,
	expires time.Time,
	err error,
) {
	return protoGetAppraisalStatus(
		x,
		r,
		code,
		include_items,
		contracts.ProtoGetHaulContractWithItems,
		contracts.ProtoGetHaulContract,
	)
}

func ProtoGetShopAppraisalStatus(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
	include_items bool,
) (
	rep ProtoAppraisalStatusRep,
	expires time.Time,
	err error,
) {
	rep.Status = proto.AppraisalStatus_AS_UNDEFINED // default

	// fetch contract status in a goroutine
	x, cancel := x.WithCancel()
	defer cancel()
	chnContract :=
		expirable.NewChanResult[*ProtoAppraisalStatusRep](x.Ctx(), 1, 0)
	if include_items {
		go expirable.P4Transceive(
			chnContract,
			x, r, code, contracts.ProtoGetShopContractWithItems,
			protoGetContractAppraisalStatusWithItems,
		)
	} else {
		go expirable.P4Transceive(
			chnContract,
			x, r, code, contracts.ProtoGetShopContract,
			protoGetContractAppraisalStatus,
		)
	}

	// check if code is in purchase queue
	var inPurchaseQueue bool
	inPurchaseQueue, expires, err = purchasequeue.InPurchaseQueue(x, code)
	if err != nil {
		return rep, expires, err
	} else if inPurchaseQueue {
		rep.Status = proto.AppraisalStatus_AS_PURCHASE_QUEUE
		return rep, expires, nil
	}

	// recv contract with items
	var contractRep *ProtoAppraisalStatusRep
	var contractExpires time.Time
	contractRep, contractExpires, err = chnContract.RecvExp()
	if err == nil && contractRep != nil {
		// overwrite purchase queue expires
		// once a code has been found in contracts, it will never be in purchase queue
		expires = contractExpires
		rep = *contractRep
	}

	return rep, expires, err
}

func protoGetContractAppraisalStatus(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
	getContract func(
		x cache.Context,
		r *protoregistry.ProtoRegistry,
		code string,
	) (
		contract *proto.Contract,
		expires time.Time,
		err error,
	),
) (
	rep *ProtoAppraisalStatusRep,
	expires time.Time,
	err error,
) {
	var contract *proto.Contract
	contract, expires, err = getContract(x, r, code)
	if err == nil && contract != nil {
		rep = &ProtoAppraisalStatusRep{
			Status:        proto.AppraisalStatus_AS_CONTRACT,
			Contract:      contract,
			ContractItems: nil,
		}
	}
	return rep, expires, err
}

func protoGetContractAppraisalStatusWithItems(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
	getContractWithItems func(
		x cache.Context,
		r *protoregistry.ProtoRegistry,
		code string,
	) (
		contractWithItems contracts.ProtoContractWithItemsRep,
		expires time.Time,
		err error,
	),
) (
	rep *ProtoAppraisalStatusRep,
	expires time.Time,
	err error,
) {
	var contractWithItems contracts.ProtoContractWithItemsRep
	contractWithItems, expires, err = getContractWithItems(x, r, code)
	if err == nil && contractWithItems.Contract != nil {
		rep = &ProtoAppraisalStatusRep{
			Status:        proto.AppraisalStatus_AS_CONTRACT,
			Contract:      contractWithItems.Contract,
			ContractItems: contractWithItems.Items,
		}
	}
	return rep, expires, err
}
