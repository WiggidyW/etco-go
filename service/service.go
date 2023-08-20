package service

import (
	"github.com/WiggidyW/weve-esi/client/contracts"
	"github.com/WiggidyW/weve-esi/client/contracts/items/multi"
	"github.com/WiggidyW/weve-esi/client/inventory"
	"github.com/WiggidyW/weve-esi/client/market/shop"
	"github.com/WiggidyW/weve-esi/client/shopqueue"
)

const (
	CORPORATION_ID    int32  = 0 // TODO
	WEB_REFRESH_TOKEN string = "TODO"
)

type Service struct {
	contractsClient          contracts.A_WC_ContractsClient
	shopQueueClient          shopqueue.A_ShopQueueClient
	multiContractItemsClient multi.MultiRateLimitingContractItemsClient
	shopPriceClient          shop.ShopPriceClient
	inventoryClient          inventory.A_InventoryClient
}
