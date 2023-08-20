package service

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/authingfwding"

	// "github.com/WiggidyW/weve-esi/client/esi/internal/raw"
	"github.com/WiggidyW/weve-esi/client/contracts"
	"github.com/WiggidyW/weve-esi/client/shopqueue"
	"github.com/WiggidyW/weve-esi/proto"
)

func (s *Service) ShopQueue(
	ctx context.Context,
	req *proto.ShopQueueRequest,
) (*proto.ShopQueueResponse, error) {
	rShopQueueRep, err := s.shopQueueClient.Fetch(
		ctx,
		authingfwding.WithAuthableParams[shopqueue.ShopQueueParams]{
			NativeRefreshToken: req.Auth.Token,
			Params: shopqueue.ShopQueueParams{
				BlockOnModify: false,
			},
		},
	)

	ok, authRep, errRep := authRepToGrpcRep(rShopQueueRep, err)
	grpcRep := &proto.ShopQueueResponse{
		Auth:  authRep,
		Error: errRep,
	}
	if !ok {
		return grpcRep, nil
	}

	var contractIds []int32
	rShopQueue := rShopQueueRep.Data.ShopQueue
	rShopContracts := rShopQueueRep.Data.ShopContracts
	if req.IncludeItems {
		grpcRep.Queue, contractIds = newPBShopQueueWithContractIds(
			rShopQueue,
			rShopContracts,
		)
	} else {
		grpcRep.Queue = newPBShopQueue(
			rShopQueue,
			rShopContracts,
		)
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
		if queueEntry.Contract != nil {
			rContractId := queueEntry.Contract.ContractId
			queueEntry.Contract.Items = rContractItems[rContractId]
		}
	}

	grpcRep.Queue.Naming = maybeFinishNamingSession(namingSession)

	return grpcRep, nil
}

func newPBShopQueueEntry(
	code string,
	contracts map[string]contracts.Contract,
) (hasContract bool, _ *proto.QueueEntry, id int32) {
	if contract, ok := contracts[code]; ok {
		// return a queueentry with a contract and the contract id
		return true, &proto.QueueEntry{
			Code:     code,
			Contract: newPBContract(contract),
		}, contract.ContractId
	} else {
		// return a queueentry with no contract
		return false, &proto.QueueEntry{
			Code: code,
		}, 0
	}
}

func newPBShopQueue(
	queueCodes []string,
	contracts map[string]contracts.Contract,
) *proto.Queue {
	queue := &proto.Queue{
		Entries: make([]*proto.QueueEntry, 0, len(queueCodes)),
	}

	for _, code := range queueCodes {
		_, entry, _ := newPBShopQueueEntry(code, contracts)
		queue.Entries = append(queue.Entries, entry)
	}

	return queue
}

func newPBShopQueueWithContractIds(
	queueCodes []string,
	contracts map[string]contracts.Contract,
) (_ *proto.Queue, ids []int32) {
	queue := &proto.Queue{
		Entries: make([]*proto.QueueEntry, 0, len(queueCodes)),
	}
	ids = make([]int32, 0, len(queueCodes))

	for _, code := range queueCodes {
		ok, entry, contractId := newPBShopQueueEntry(code, contracts)
		queue.Entries = append(queue.Entries, entry)
		if ok {
			ids = append(ids, contractId)
		}
	}

	return queue, ids
}
