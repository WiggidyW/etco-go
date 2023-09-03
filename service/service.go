package service

import (
	"github.com/WiggidyW/etco-go/client/auth"
	"github.com/WiggidyW/etco-go/client/bucket"
	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/client/purchase"
	rdbc "github.com/WiggidyW/etco-go/client/remotedb"
	"github.com/WiggidyW/etco-go/client/userdata"
	"github.com/WiggidyW/etco-go/staticdb"
)

// TODO: check if auth is null in requests
type Service struct {
	authClient                     auth.AuthClient
	rReadAuthListClient            bucket.AuthListReaderClient
	rWriteAuthListClient           bucket.AuthListWriterClient
	rUserDataClient                userdata.UserDataClient
	rdbcUserDataClient             rdbc.SC_ReadUserDataClient
	shopInventoryClient            protoclient.PBShopInventoryClient
	shopContractQueueClient        protoclient.PBShopContractQueueClient
	buybackContractQueueClient     protoclient.PBBuybackContractQueueClient
	shopPurchaseQueueClient        protoclient.PBShopPurchaseQueueClient
	shopMakePurchaseClient         purchase.MakePurchaseClient
	shopCancelPurchaseClient       purchase.CancelPurchaseClient
	shopDeletePurchasesClient      rdbc.SMAC_DelPurchasesClient
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
	shopLocationsClient            protoclient.ShopLocationsClient
}
