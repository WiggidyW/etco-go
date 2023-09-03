package inventory

import (
	"context"

	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go/client/inventory/locationassets"
	"github.com/WiggidyW/etco-go/client/shopqueue"
)

type InventoryClient struct {
	shopQueueClient shopqueue.ShopQueueClient
	assetsClient    locationassets.LocationShopAssetsClient
}

func NewInventoryClient(
	shopQueueClient shopqueue.ShopQueueClient,
	locationAssetsClient locationassets.LocationShopAssetsClient,
) InventoryClient {
	return InventoryClient{shopQueueClient, locationAssetsClient}
}

func (bic InventoryClient) Fetch(
	ctx context.Context,
	params InventoryParams,
) (*map[int32]int64, error) {
	chnModifiedSend, chnModifiedRecv := chanresult.
		NewChanResult[struct{}](ctx, 0, 0).Split()

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
		locationassets.LocationShopAssetsParams{
			ShopQueue:  sqRep.ParsedShopQueue,
			LocationId: params.LocationId,
		},
	)
	if err != nil {
		return nil, err
	}

	return &inventory, nil
}
