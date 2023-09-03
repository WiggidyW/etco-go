package contracts

import (
	"context"

	"github.com/WiggidyW/chanresult"
)

type MultiContractItemsClient struct {
	Inner WC_SingleContractItemsClient
}

func (mrlcic MultiContractItemsClient) Fetch(
	ctx context.Context,
	params MultiContractItemsParams,
) (chanresult.ChanRecvResult[MultiContractItems], error) {
	chnSend, chnRecv := chanresult.
		NewChanResult[MultiContractItems](ctx, 0, 0).Split()
	go mrlcic.fetchAll(ctx, params, chnSend)
	return chnRecv, nil
}

// sends out requests for all contract ids, ensuring that no more than
// i.MAX_CONCURRENT_REQS are active at any given time
func (mrlcic MultiContractItemsClient) fetchAll(
	ctx context.Context,
	params MultiContractItemsParams,
	chnSendRep chanresult.ChanSendResult[MultiContractItems],
) {
	var activeReqs int = 0
	chnSendDone := make(chan struct{}, MAX_CONCURRENT_REQS)

	for _, contractId := range params.ContractIds {
		if activeReqs >= MAX_CONCURRENT_REQS {
			select {
			case <-ctx.Done():
				return
			case <-chnSendDone:
				activeReqs--
			}
		}
		go mrlcic.fetchSingle(
			ctx,
			SingleContractItemsParams{
				ContractId: contractId,
			},
			chnSendRep,
			chnSendDone,
		)
		activeReqs++
	}
}

// fetches a single contracts items and sends the result to chnSendRep
// also sends a struct{} to chnSendDone, signifying that the request is done
func (mrlcic MultiContractItemsClient) fetchSingle(
	ctx context.Context,
	params SingleContractItemsParams,
	chnSendRep chanresult.ChanSendResult[MultiContractItems],
	chnSendDone chan<- struct{},
) {
	rep, err := mrlcic.Inner.Fetch(ctx, params)
	chnSendDone <- struct{}{}
	if err != nil {
		chnSendRep.SendErr(err)
	} else {
		chnSendRep.SendOk(MultiContractItems{
			ContractId:    params.ContractId,
			ContractItems: rep.Data(),
		})
	}
}
