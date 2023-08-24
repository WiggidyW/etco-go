package service

import (
	"github.com/WiggidyW/eve-trading-co-go/client/appraisal/buyback"
	shopappr "github.com/WiggidyW/eve-trading-co-go/client/appraisal/shop"
	admin "github.com/WiggidyW/eve-trading-co-go/client/configure/authlist"
	getbuybackbuilder "github.com/WiggidyW/eve-trading-co-go/client/configure/btypemapsbuilder/get"
	mergebuybackbuilder "github.com/WiggidyW/eve-trading-co-go/client/configure/btypemapsbuilder/pbmerge"
	getbuybacksystems "github.com/WiggidyW/eve-trading-co-go/client/configure/buybacksystems/get"
	mergebuybacksystems "github.com/WiggidyW/eve-trading-co-go/client/configure/buybacksystems/pbmerge"
	getmarkets "github.com/WiggidyW/eve-trading-co-go/client/configure/markets/get"
	mergemarkets "github.com/WiggidyW/eve-trading-co-go/client/configure/markets/pbmerge"
	getshoplocations "github.com/WiggidyW/eve-trading-co-go/client/configure/shoplocations/get"
	mergeshoplocations "github.com/WiggidyW/eve-trading-co-go/client/configure/shoplocations/pbmerge"
	getshopbuilder "github.com/WiggidyW/eve-trading-co-go/client/configure/stypemapsbuilder/get"
	mergeshopbuilder "github.com/WiggidyW/eve-trading-co-go/client/configure/stypemapsbuilder/pbmerge"
	"github.com/WiggidyW/eve-trading-co-go/client/contracts"
	"github.com/WiggidyW/eve-trading-co-go/client/contracts/items/multi"
	"github.com/WiggidyW/eve-trading-co-go/client/inventory"
	shopmarket "github.com/WiggidyW/eve-trading-co-go/client/market/shop"
	"github.com/WiggidyW/eve-trading-co-go/client/shopqueue"
)

const (
	CORPORATION_ID    int32  = 0 // TODO
	WEB_REFRESH_TOKEN string = "TODO"
)

// TODO: check if auth is null in requests
type Service struct {
	contractsClient           contracts.A_WC_ContractsClient
	shopQueueClient           shopqueue.A_ShopQueueClient
	multiContractItemsClient  multi.MultiRateLimitingContractItemsClient
	shopPriceClient           shopmarket.ShopPriceClient
	inventoryClient           inventory.A_InventoryClient
	anonBuybackApprClient     buyback.BuybackAppraisalClient
	charBuybackApprClient     buyback.F_BuybackAppraisalClient
	admnShopApprClient        shopappr.AF_ShopAppraisalClient
	getShopLocationsClient    getshoplocations.A_GetShopLocationsClient
	mergeShopLocationsClient  mergeshoplocations.A_PbMergeShopLocationsClient
	getBuybackSystemsClient   getbuybacksystems.A_GetBuybackSystemsClient
	mergeBuybackSystemsClient mergebuybacksystems.A_PbMergeBuybackSystemsBuilderClient
	getMarketsClient          getmarkets.A_GetMarketsClient
	mergeMarketsClient        mergemarkets.A_PbMergeMarketsClient
	getBuybackBuilderClient   getbuybackbuilder.A_GetBuybackSystemTypeMapsBuilderClient
	mergeBuybackBuilderClient mergebuybackbuilder.A_PbMergeBuybackSystemTypeMapsBuilderClient
	getShopBuilderClient      getshopbuilder.A_GetShopLocationTypeMapsBuilderClient
	mergeShopBuilderClient    mergeshopbuilder.A_PbMergeShopLocationTypeMapsBuilderClient
	getAuthListClient         admin.A_AdminReadClient
	setAuthListClient         admin.A_AdminWriteClient
}
