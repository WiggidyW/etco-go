package service

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/authingfwding"
	"github.com/WiggidyW/weve-esi/client/contracts"
	"github.com/WiggidyW/weve-esi/proto"
)

func (s *Service) BuybackQueue(
	ctx context.Context,
	req *proto.BuybackQueueRequest,
) (*proto.BuybackQueueResponse, error) {
	rContractsRep, err := s.contractsClient.Fetch(
		ctx,
		authingfwding.WithAuthableParams[contracts.ContractsParams]{
			NativeRefreshToken: req.Auth.Token,
			Params:             contracts.ContractsParams{},
		},
	)

	ok, authRep, errRep := authRepToGrpcRep(rContractsRep, err)
	grpcRep := &proto.BuybackQueueResponse{
		Auth:  authRep,
		Error: errRep,
	}
	if !ok {
		return grpcRep, nil
	}

	var contractIds []int32
	rBuybackContracts := rContractsRep.Data.Data().BuybackContracts
	if req.IncludeItems {
		grpcRep.Queue, contractIds = newPBBuybackQueueWithContractIds(
			rBuybackContracts,
		)
	} else {
		grpcRep.Queue = newPBBuybackQueue(rBuybackContracts)
		return grpcRep, nil
	}

	namingSession := maybeNewLocalNamingSession(req.IncludeNaming)
	rContractItems, err := s.fetchAllContractItems(
		ctx,
		contractIds,
		namingSession,
	)
	if err != nil {
		grpcRep.Error = newErrorResponse(err)
		grpcRep.Queue = nil
		return grpcRep, nil
	}

	for _, queueEntry := range grpcRep.Queue.Entries {
		rContractId := queueEntry.Contract.ContractId
		queueEntry.Contract.Items = rContractItems[rContractId]
	}

	grpcRep.Queue.Naming = maybeFinishNamingSession(namingSession)

	return grpcRep, nil
}

func newPBBuybackQueueEntry(
	code string,
	contract contracts.Contract,
) (_ *proto.QueueEntry, id int32) {
	return &proto.QueueEntry{
		Code:     code,
		Contract: newPBContract(contract),
	}, contract.ContractId
}

func newPBBuybackQueue(
	contracts map[string]contracts.Contract,
) *proto.Queue {
	queue := &proto.Queue{
		Entries: make([]*proto.QueueEntry, 0, len(contracts)),
	}

	for code, contract := range contracts {
		entry, _ := newPBBuybackQueueEntry(code, contract)
		queue.Entries = append(queue.Entries, entry)
	}

	return queue
}

func newPBBuybackQueueWithContractIds(
	contracts map[string]contracts.Contract,
) (_ *proto.Queue, ids []int32) {
	queue := &proto.Queue{
		Entries: make([]*proto.QueueEntry, 0, len(contracts)),
	}
	ids = make([]int32, 0, len(contracts))

	for code, contract := range contracts {
		entry, id := newPBBuybackQueueEntry(code, contract)
		queue.Entries = append(queue.Entries, entry)
		ids = append(ids, id)
	}

	return queue, ids
}
