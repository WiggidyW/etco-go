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
	ParsedShopQueue,
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
) (*ParsedShopQueue, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the shop contracts in a separate goroutine
	chnSend, chnRecv := util.NewChanResult[map[string]contracts.Contract](
		ctx,
	).Split()
	go sqc.fetchContracts(ctx, contracts.ContractsParams{}, chnSend)

	// fetch the raw shop queue
	readRep, err := sqc.readClient.Fetch(ctx, read.ShopQueueReadParams{})
	if err != nil {
		return nil, err
	}
	readQueue := readRep.Data()

	// wait for the shop contracts
	okQueue := make([]string, 0, len(readQueue))
	delQueue := make([]string, 0, len(readQueue))
	contracts, err := chnRecv.Recv()
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

	// if the delete queue has entries, remove them
	if modified {
		if params.BlockOnModify {
			if err := sqc.deleteCodes(ctx, delQueue); err != nil {
				return nil, err
			}
		} else {
			go func() {
				if err := sqc.deleteCodes(
					context.Background(),
					delQueue,
				); err != nil {
					logger.Err(err)
				}
			}()
		}
	}

	return &ParsedShopQueue{
		ShopQueue:     okQueue,
		ShopContracts: contracts,
		Modified:      modified,
	}, nil
}

func (sqc ShopQueueClient) deleteCodes(
	ctx context.Context,
	delQueue []string,
) error {
	_, err := sqc.removeClient.Fetch(
		ctx,
		removematching.ShopQueueRemoveMatchingParams(
			delQueue,
		),
	)
	return err
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
