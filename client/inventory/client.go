package inventory

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/authingfwding"
	"github.com/WiggidyW/weve-esi/client/authingfwding/authing"
	"github.com/WiggidyW/weve-esi/client/inventory/internal/location"
	"github.com/WiggidyW/weve-esi/client/shopqueue"
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
	sqRep, err := bic.shopQueueClient.Fetch(
		ctx,
		shopqueue.NewShopQueueParams(
			true, // block on modify to avoid cache racing
		),
	)
	if err != nil {
		return nil, err
	}

	inventory, err := bic.assetsClient.Fetch(
		ctx,
		location.NewLocationShopAssetsParams(
			sqRep.ShopQueue,
			params.LocationId,
		),
	)
	if err != nil {
		return nil, err
	}

	return &inventory, nil
}
