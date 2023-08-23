package inventory

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/authingfwding"
	"github.com/WiggidyW/weve-esi/client/authingfwding/authing"
	"github.com/WiggidyW/weve-esi/client/inventory/internal/location"
	"github.com/WiggidyW/weve-esi/client/shopqueue"
	"github.com/WiggidyW/weve-esi/util"
)

type A_InventoryClient = authing.AuthingClient[
	authingfwding.WithAuthableParams[InventoryParams],
	InventoryParams,
	map[int32]int64,
	InventoryClient,
]

type InventoryClient struct {
	shopQueueClient shopqueue.ShopQueueClient
	assetsClient    location.LocationShopAssetsClient
}

func (bic InventoryClient) Fetch(
	ctx context.Context,
	params InventoryParams,
) (*map[int32]int64, error) {
	chnModified := util.NewChanResult[struct{}](ctx)
	chnModifiedSend, chnModifiedRecv := chnModified.Split()

	sqRep, err := bic.shopQueueClient.Fetch(
		ctx,
		// block on modify to avoid cache racing
		shopqueue.ShopQueueParams{ChnSendModifyDone: &chnModifiedSend},
	)
	if err != nil {
		return nil, err
	}

	if params.ChnSendShopQueueRep != nil { // if a channel was provided
		// send the shop queue, but don't block
		go func() {
			// if context was cancelled, we'll find out soon enough
			_ = params.ChnSendShopQueueRep.SendOk(sqRep)
		}()
	}

	if sqRep.Modified {
		// wait for modification to finish
		_, err = chnModifiedRecv.Recv()
		if err != nil {
			return nil, err
		}
	}

	inventory, err := bic.assetsClient.Fetch(
		ctx,
		location.LocationShopAssetsParams{
			ShopQueue:  sqRep.ParsedShopQueue,
			LocationId: params.LocationId,
		},
	)
	if err != nil {
		return nil, err
	}

	return &inventory, nil
}
