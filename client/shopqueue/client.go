package shopqueue

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/authingfwding"
	"github.com/WiggidyW/weve-esi/client/authingfwding/authing"
	"github.com/WiggidyW/weve-esi/client/contracts"
	"github.com/WiggidyW/weve-esi/client/remotedb/rawshopqueue/read"
	"github.com/WiggidyW/weve-esi/client/remotedb/rawshopqueue/removematching"
	"github.com/WiggidyW/weve-esi/logger"
	"github.com/WiggidyW/weve-esi/util"
)

type A_ShopQueueClient = authing.AuthingClient[
	authingfwding.WithAuthableParams[ShopQueueParams],
	ShopQueueParams,
	ShopQueueResponse,
	ShopQueueClient,
]

type ShopQueueClient struct {
	readClient      read.SC_ShopQueueReadClient
	removeClient    removematching.SMAC_ShopQueueRemoveMatchingClient
	contractsClient contracts.WC_ContractsClient
}

// returns a shop queue that only includes codes that do not yet have an ESI contract
// also returns the contracts
func (sqc ShopQueueClient) Fetch(
	ctx context.Context,
	params ShopQueueParams,
) (*ShopQueueResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the shop contracts in a separate goroutine
	chnContracts := util.NewChanResult[map[string]contracts.Contract](ctx)
	chnSendContracts, chnRecvContracts := chnContracts.Split()
	go sqc.fetchContracts(ctx, contracts.ContractsParams{}, chnSendContracts)

	// fetch the raw shop queue
	readRep, err := sqc.readClient.Fetch(ctx, read.ShopQueueReadParams{})
	if err != nil {
		return nil, err
	}
	readQueue := readRep.Data()

	// wait for the shop contracts
	okQueue := make([]string, 0, len(readQueue))
	delQueue := make([]string, 0, len(readQueue))
	contracts, err := chnRecvContracts.Recv()
	if err != nil {
		return nil, err
	}

	// filter the shop queue
	// TODO: Make a new client that just gets a hashset of the shop contracts
	for _, code := range readQueue {
		if _, ok := contracts[code]; !ok {
			okQueue = append(okQueue, code)
		} else {
			delQueue = append(delQueue, code)
		}
	}
	modified := len(delQueue) > 0

	if modified {
		go sqc.handleModify(params.ChnSendModifyDone, delQueue)
	}

	return &ShopQueueResponse{
		ParsedShopQueue: okQueue,
		ShopContracts:   contracts,
		Modified:        modified,
	}, nil
}

func (sqc ShopQueueClient) fetchContracts(
	ctx context.Context,
	params contracts.ContractsParams,
	chnSend util.ChanSendResult[map[string]contracts.Contract],
) {
	if contractsRep, err := sqc.contractsClient.Fetch(
		ctx,
		params,
	); err != nil {
		chnSend.SendErr(err)
	} else {
		chnSend.SendOk(contractsRep.Data().ShopContracts)
	}
}

func (sqc ShopQueueClient) handleModify(
	chnSendModifyDone *util.ChanSendResult[struct{}],
	delQueue []string,
) error {
	_, err := sqc.removeClient.Fetch(
		context.Background(),
		removematching.ShopQueueRemoveMatchingParams(delQueue),
	)
	if err != nil {
		return sendModifyResult(chnSendModifyDone, err)
	} else {
		return sendModifyResult(chnSendModifyDone, nil)
	}
}

// - error is not nil, channel is     nil, logs error
// - error is not nil, channel is not nil, sends error
// - error is nil,     channel is     nil, does nothing
// - error is nil,     channel is not nil, sends struct{}
func sendModifyResult(
	chnSendModifyDone *util.ChanSendResult[struct{}],
	err error,
) error /* ctx error */ {
	if err != nil {
		if chnSendModifyDone == nil {
			logger.Err(err)
			return nil
		} else {
			return chnSendModifyDone.SendErr(err)
		}
	} else if chnSendModifyDone == nil {
		return nil
	} else {
		return chnSendModifyDone.SendOk(struct{}{})
	}
}
