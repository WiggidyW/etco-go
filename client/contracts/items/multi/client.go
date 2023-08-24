package multi

import (
	"context"

	i "github.com/WiggidyW/eve-trading-co-go/client/contracts/items"
	"github.com/WiggidyW/eve-trading-co-go/client/contracts/items/single"
	"github.com/WiggidyW/eve-trading-co-go/util"
)

type MultiRateLimitingContractItemsClient struct {
	Inner single.WC_RateLimitingContractItemsClient
}

func (mrlcic MultiRateLimitingContractItemsClient) Fetch(
	ctx context.Context,
	params MultiRateLimitingContractItemsParams,
) (util.ChanRecvResult[ContractItems], error) {
	chnSend, chnRecv := util.NewChanResult[ContractItems](ctx).Split()
	go mrlcic.fetchAll(ctx, params, chnSend)
	return chnRecv, nil
}

// sends out requests for all contract ids, ensuring that no more than
// i.MAX_CONCURRENT_REQS are active at any given time
func (mrlcic MultiRateLimitingContractItemsClient) fetchAll(
	ctx context.Context,
	params MultiRateLimitingContractItemsParams,
	chnSendRep util.ChanSendResult[ContractItems],
) {
	var activeReqs int = 0
	chnSendDone := make(chan struct{}, i.MAX_CONCURRENT_REQS)

	for _, contractId := range params.ContractIds {
		if activeReqs >= i.MAX_CONCURRENT_REQS {
			select {
			case <-ctx.Done():
				return
			case <-chnSendDone:
				activeReqs--
			}
		}
		go mrlcic.fetchOne(
			ctx,
			single.RateLimitingContractItemsParams{
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
func (mrlcic MultiRateLimitingContractItemsClient) fetchOne(
	ctx context.Context,
	params single.RateLimitingContractItemsParams,
	chnSendRep util.ChanSendResult[ContractItems],
	chnSendDone chan<- struct{},
) {
	rep, err := mrlcic.Inner.Fetch(ctx, params)
	chnSendDone <- struct{}{}
	if err != nil {
		chnSendRep.SendErr(err)
	} else {
		chnSendRep.SendOk(ContractItems{
			ContractId:    params.ContractId,
			ContractItems: rep.Data(),
		})
	}
}
