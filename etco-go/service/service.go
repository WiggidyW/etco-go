package service

import (
	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/staticdb"
)

// TODO: check if auth is null in requests
type Service struct {
	shopInventoryClient            protoclient.PBShopInventoryClient
	shopContractQueueClient        protoclient.PBShopContractQueueClient
	buybackContractQueueClient     protoclient.PBBuybackContractQueueClient
	shopPurchaseQueueClient        protoclient.PBShopPurchaseQueueClient
	newBuybackAppraisalClient      protoclient.PBNewBuybackAppraisalClient[*staticdb.LocalIndexMap]
	newShopAppraisalClient         protoclient.PBNewShopAppraisalClient[*staticdb.LocalIndexMap]
	getBuybackAppraisalClient      protoclient.PBGetBuybackAppraisalClient[*staticdb.LocalIndexMap]
	getShopAppraisalClient         protoclient.PBGetShopAppraisalClient[*staticdb.LocalIndexMap]
	statusBuybackAppraisalClient   protoclient.PBStatusBuybackAppraisalClient
	statusShopAppraisalClient      protoclient.PBStatusShopAppraisalClient
	cfgMergeBTypeMapsBuilderClient protoclient.CfgMergeBuybackSystemTypeMapsBuilderClient
	cfgGetBTypeMapsBuilderClient   protoclient.CfgGetBuybackSystemTypeMapsBuilderClient
	cfgMergeSTypeMapsBuilderClient protoclient.CfgMergeShopLocationTypeMapsBuilderClient
	cfgGetSTypeMapsBuilderClient   protoclient.CfgGetShopLocationTypeMapsBuilderClient
	cfgMergeBuybackSystemsClient   protoclient.CfgMergeBuybackSystemsClient
	cfgGetBuybackSystemsClient     protoclient.CfgGetBuybackSystemsClient
	cfgMergeShopLocationsClient    protoclient.CfgMergeShopLocationsClient
	cfgGetShopLocationsClient      protoclient.CfgGetShopLocationsClient
	cfgMergeMarketsClient          protoclient.CfgMergeMarketsClient
	cfgGetMarketsClient            protoclient.CfgGetMarketsClient
	cfgGetBuybackBundleKeysClient  protoclient.CfgGetBuybackBundleKeysClient
	cfgGetShopBundleKeysClient     protoclient.CfgGetShopBundleKeysClient
	cfgGetMarketNamesClient        protoclient.CfgGetMarketNamesClient
	shopLocationsClient            protoclient.PBShopLocationsClient
	proto.UnimplementedEveTradingCoServer
}

func NewService() *Service {
	// Proto Composition clients (Non-CFG)

	localPBContractItemsClient := protoclient.
		NewPBContractItemsClient[*staticdb.LocalIndexMap]()
	localPBGetBuybackAppraisalClient := protoclient.
		NewPBGetBuybackAppraisalClient[*staticdb.LocalIndexMap]()
	localPBGetShopAppraisalClient := protoclient.
		NewPBGetShopAppraisalClient[*staticdb.LocalIndexMap]()
	localPBNewBuybackAppraisalClient := protoclient.
		NewPBNewBuybackAppraisalClient[*staticdb.LocalIndexMap]()
	localPBNewShopAppraisalClient := protoclient.
		NewPBNewShopAppraisalClient[*staticdb.LocalIndexMap]()
	statusBuybackAppraisalClient := protoclient.
		NewPBStatusBuybackAppraisalClient(
			localPBContractItemsClient,
		)
	statusShopAppraisalClient := protoclient.NewPBStatusShopAppraisalClient(
		localPBContractItemsClient,
	)
	pbBuybackContractQueueClient := protoclient.
		NewPBBuybackContractQueueClient()
	pbShopContractQueueClient := protoclient.NewPBShopContractQueueClient()
	pbShopPurchaseQueueClient := protoclient.NewPBShopPurchaseQueueClient()
	pbShopInventoryClient := protoclient.NewPBShopInventoryClient()
	pbShopLocationsClient := protoclient.NewPBShopLocationsClient()

	// Proto Composition clients (CFG)

	cfgGetBTypeMapsBuilderClient := protoclient.
		NewCfgGetBuybackSystemTypeMapsBuilderClient()
	cfgMergeBTypeMapsBuilderClient := protoclient.
		NewCfgMergeBuybackSystemTypeMapsBuilderClient()
	cfgGetSTypeMapsBuilderClient := protoclient.
		NewCfgGetShopLocationTypeMapsBuilderClient()
	cfgMergeSTypeMapsBuilderClient := protoclient.
		NewCfgMergeShopLocationTypeMapsBuilderClient()
	cfgGetBuybackSystemsClient := protoclient.NewCfgGetBuybackSystemsClient()
	cfgMergeBuybackSystemsClient := protoclient.
		NewCfgMergeBuybackSystemsClient()
	cfgGetShopLocationsClient := protoclient.NewCfgGetShopLocationsClient()
	cfgMergeShopLocationsClient := protoclient.
		NewCfgMergeShopLocationsClient()
	cfgGetMarketsClient := protoclient.NewCfgGetMarketsClient()
	cfgMergeMarketsClient := protoclient.NewCfgMergeMarketsClient()
	cfgGetBuybackBundleKeysClient := protoclient.NewCfgGetBuybackBundleKeysClient()
	cfgGetShopBundleKeysClient := protoclient.NewCfgGetShopBundleKeysClient()
	cfgGetMarketNamesClient := protoclient.NewCfgGetMarketNamesClient()

	return &Service{
		shopInventoryClient:            pbShopInventoryClient,
		shopContractQueueClient:        pbShopContractQueueClient,
		buybackContractQueueClient:     pbBuybackContractQueueClient,
		shopPurchaseQueueClient:        pbShopPurchaseQueueClient,
		newBuybackAppraisalClient:      localPBNewBuybackAppraisalClient,
		newShopAppraisalClient:         localPBNewShopAppraisalClient,
		getBuybackAppraisalClient:      localPBGetBuybackAppraisalClient,
		getShopAppraisalClient:         localPBGetShopAppraisalClient,
		statusBuybackAppraisalClient:   statusBuybackAppraisalClient,
		statusShopAppraisalClient:      statusShopAppraisalClient,
		cfgMergeBTypeMapsBuilderClient: cfgMergeBTypeMapsBuilderClient,
		cfgGetBTypeMapsBuilderClient:   cfgGetBTypeMapsBuilderClient,
		cfgMergeSTypeMapsBuilderClient: cfgMergeSTypeMapsBuilderClient,
		cfgGetSTypeMapsBuilderClient:   cfgGetSTypeMapsBuilderClient,
		cfgMergeBuybackSystemsClient:   cfgMergeBuybackSystemsClient,
		cfgGetBuybackSystemsClient:     cfgGetBuybackSystemsClient,
		cfgMergeShopLocationsClient:    cfgMergeShopLocationsClient,
		cfgGetShopLocationsClient:      cfgGetShopLocationsClient,
		cfgMergeMarketsClient:          cfgMergeMarketsClient,
		cfgGetMarketsClient:            cfgGetMarketsClient,
		cfgGetBuybackBundleKeysClient:  cfgGetBuybackBundleKeysClient,
		cfgGetShopBundleKeysClient:     cfgGetShopBundleKeysClient,
		cfgGetMarketNamesClient:        cfgGetMarketNamesClient,
		shopLocationsClient:            pbShopLocationsClient,
	}
}
