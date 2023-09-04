package service

import (
	"net/http"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/appraisal"
	"github.com/WiggidyW/etco-go/client/auth"
	bucketc "github.com/WiggidyW/etco-go/client/bucket"
	"github.com/WiggidyW/etco-go/client/contracts"
	"github.com/WiggidyW/etco-go/client/esi/jwt"
	massetscorporation "github.com/WiggidyW/etco-go/client/esi/model/assetscorporation"
	mcharacterinfo "github.com/WiggidyW/etco-go/client/esi/model/characterinfo"
	mcontractitems "github.com/WiggidyW/etco-go/client/esi/model/contractitems"
	mcontractscorporation "github.com/WiggidyW/etco-go/client/esi/model/contractscorporation"
	mordersregion "github.com/WiggidyW/etco-go/client/esi/model/ordersregion"
	mordersstructure "github.com/WiggidyW/etco-go/client/esi/model/ordersstructure"
	mstructureinfo "github.com/WiggidyW/etco-go/client/esi/model/structureinfo"
	"github.com/WiggidyW/etco-go/client/esi/raw_"
	"github.com/WiggidyW/etco-go/client/inventory"
	"github.com/WiggidyW/etco-go/client/inventory/locationassets"
	"github.com/WiggidyW/etco-go/client/market"
	"github.com/WiggidyW/etco-go/client/market/marketprice"
	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/client/purchase"
	rdbc "github.com/WiggidyW/etco-go/client/remotedb"
	"github.com/WiggidyW/etco-go/client/shopqueue"
	"github.com/WiggidyW/etco-go/client/structureinfo"
	"github.com/WiggidyW/etco-go/client/userdata"
	"github.com/WiggidyW/etco-go/proto"
	rdb "github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

// TODO: check if auth is null in requests
type Service struct {
	authClient                     auth.AuthClient
	rReadAuthListClient            bucketc.AuthListReaderClient
	rWriteAuthListClient           bucketc.AuthListWriterClient
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
	shopLocationsClient            protoclient.PBShopLocationsClient
	proto.UnimplementedEveTradingCoServer
}

func NewService() *Service {
	// basal clients

	sharedServerCache := cache.NewSharedServerCache()
	sharedClientCache := cache.NewSharedClientCache()
	rawBucketClient := bucket.NewBucketClient()
	rawRemoteDBClient := rdb.NewRemoteDBClient()
	httpClient := &http.Client{}

	// raw ESI clients

	unauthRawClient := raw_.NewUnauthenticatedRawClient(httpClient)
	corpRawClient := raw_.NewCorpRawClient(httpClient)
	marketsRawClient := raw_.NewMarketsRawClient(httpClient)
	structureInfoRawClient := raw_.NewStructureInfoRawClient(httpClient)
	authRawClient := raw_.NewAuthRawClient(httpClient)

	// Higher Level ESI clients (Model clients + JWT client)

	mContractsCorporationClient :=
		mcontractscorporation.NewContractsCorporationClient(corpRawClient)
	mAssetsCorporationClient :=
		massetscorporation.NewAssetsCorporationClient(corpRawClient)
	mContractItemsClient :=
		mcontractitems.NewContractItemsClient(corpRawClient)
	mOrdersStructureClient :=
		mordersstructure.NewOrdersStructureClient(marketsRawClient)
	mStructureInfoClient :=
		mstructureinfo.NewStructureInfoClient(structureInfoRawClient)
	mOrdersRegionClient :=
		mordersregion.NewOrdersRegionClient(unauthRawClient)
	wc_mCharacterInfoClient := mcharacterinfo.NewWC_CharacterInfoClient(
		unauthRawClient,
		sharedClientCache,
		sharedServerCache,
	)
	jwtClient := jwt.NewJWTClient(
		unauthRawClient,
		authRawClient,
		sharedClientCache,
		sharedServerCache,
	)

	// Higher level bucket clients

	authHashSetReaderClient := bucketc.NewSC_AuthHashSetReaderClient(
		rawBucketClient,
		sharedServerCache,
	)
	authHashSetWriterClient := bucketc.NewSAC_AuthHashSetWriterClient(
		rawBucketClient,
		authHashSetReaderClient.GetAntiCache(),
	)
	authListReaderClient := bucketc.NewAuthListReaderClient(
		authHashSetReaderClient,
	)
	authListWriterClient := bucketc.NewAuthListWriterClient(
		authHashSetWriterClient,
	)
	webBTypeMapsBuilderReaderClient :=
		bucketc.NewSC_WebBuybackSystemTypeMapsBuilderReaderClient(
			rawBucketClient,
			sharedServerCache,
		)
	webBTypeMapsBuilderWriterClient :=
		bucketc.NewSAC_WebBuybackSystemTypeMapsBuilderWriterClient(
			rawBucketClient,
			webBTypeMapsBuilderReaderClient.GetAntiCache(),
		)
	webSTypeMapsBuilderReaderClient :=
		bucketc.NewSC_WebShopLocationTypeMapsBuilderReaderClient(
			rawBucketClient,
			sharedServerCache,
		)
	webSTypeMapsBuilderWriterClient :=
		bucketc.NewSAC_WebShopLocationTypeMapsBuilderWriterClient(
			rawBucketClient,
			webSTypeMapsBuilderReaderClient.GetAntiCache(),
		)
	webBuybackSystemsReaderClient :=
		bucketc.NewSC_WebBuybackSystemsReaderClient(
			rawBucketClient,
			sharedServerCache,
		)
	webBuybackSystemsWriterClient :=
		bucketc.NewSAC_WebBuybackSystemsWriterClient(
			rawBucketClient,
			webBuybackSystemsReaderClient.GetAntiCache(),
		)
	webShopLocationsReaderClient :=
		bucketc.NewSC_WebShopLocationsReaderClient(
			rawBucketClient,
			sharedServerCache,
		)
	webShopLocationsWriterClient :=
		bucketc.NewSAC_WebShopLocationsWriterClient(
			rawBucketClient,
			webShopLocationsReaderClient.GetAntiCache(),
		)
	webMarketsReaderClient := bucketc.NewSC_WebMarketsReaderClient(
		rawBucketClient,
		sharedServerCache,
	)
	webMarketsWriterClient := bucketc.NewSAC_WebMarketsWriterClient(
		rawBucketClient,
		webMarketsReaderClient.GetAntiCache(),
	)

	// Higher Level remoteDB clients + Unreserved Location Assets

	wc_rdbcReadShopAppraisalClient := rdbc.NewWC_ReadShopAppraisalClient(
		rawRemoteDBClient,
		sharedClientCache,
		sharedServerCache,
	)
	wc_rdbcReadBuybackAppraisalClient :=
		rdbc.NewWC_ReadBuybackAppraisalClient(
			rawRemoteDBClient,
			sharedClientCache,
			sharedServerCache,
		)
	sc_rdbcReadShopQueueClient := rdbc.NewSC_ReadShopQueueClient(
		rawRemoteDBClient,
		sharedServerCache,
	)
	sc_rdbcReadUserDataClient := rdbc.NewSC_ReadUserDataClient(
		rawRemoteDBClient,
		sharedServerCache,
	)
	locationShopAssetsClient := locationassets.NewLocationShopAssetsClient(
		mAssetsCorporationClient,
		wc_rdbcReadShopAppraisalClient,
		sharedClientCache,
		sharedServerCache,
	)

	rdbcReadShopQueueAntiCache := sc_rdbcReadShopQueueClient.GetAntiCache()
	rdbcReadUserDataAntiCache := sc_rdbcReadUserDataClient.GetAntiCache()
	unreservedAntiCache := locationShopAssetsClient.GetUnreservedAntiCache()

	sac_rdbcWriteBuybackAppraisalClient :=
		rdbc.NewSAC_WriteBuybackAppraisalClient(
			rawRemoteDBClient,
			rdbcReadUserDataAntiCache,
		)
	smac_rdbcDelPurchasesClient := rdbc.NewSMAC_DelPurchasesClient(
		rawRemoteDBClient,
		rdbcReadShopQueueAntiCache,
		unreservedAntiCache,
	)
	smac_rdbcCancelPurchaseClient := rdbc.NewSMAC_CancelPurchaseClient(
		rawRemoteDBClient,
		rdbcReadUserDataAntiCache,
		rdbcReadShopQueueAntiCache,
		unreservedAntiCache,
	)
	smac_rdbcWritePurchaseClient := rdbc.NewSMAC_WritePurchaseClient(
		rawRemoteDBClient,
		rdbcReadUserDataAntiCache,
		rdbcReadShopQueueAntiCache,
		unreservedAntiCache,
	)

	// Non-Proto Composition clients

	marketPriceClient := marketprice.NewMarketPriceClient(
		mOrdersRegionClient,
		mOrdersStructureClient,
		sharedClientCache,
		sharedServerCache,
	)
	buybackPriceClient := market.NewBuybackPriceClient(marketPriceClient)
	shopPriceClient := market.NewShopPriceClient(marketPriceClient)
	makeBuybackAppraisalClient := appraisal.NewMakeBuybackAppraisalClient(
		sac_rdbcWriteBuybackAppraisalClient,
		buybackPriceClient,
	)
	makeShopAppraisalClient := appraisal.NewMakeShopAppraisalClient(
		shopPriceClient,
	)
	wc_StructureInfoClient := structureinfo.NewWC_StructureInfoClient(
		mStructureInfoClient,
		sharedClientCache,
		sharedServerCache,
	)
	wc_ContractsClient := contracts.NewWC_ContractsClient(
		mContractsCorporationClient,
		sharedClientCache,
		sharedServerCache,
	)
	wc_ContractItemsClient := contracts.NewWC_SingleContractItemsClient(
		mContractItemsClient,
		sharedClientCache,
		sharedServerCache,
	)
	shopQueueClient := shopqueue.NewShopQueueClient(
		sc_rdbcReadShopQueueClient,
		smac_rdbcDelPurchasesClient,
		wc_ContractsClient,
	)
	shopInventoryClient := inventory.NewInventoryClient(
		shopQueueClient,
		locationShopAssetsClient,
	)
	userDataClient := userdata.NewUserDataClient(
		sc_rdbcReadUserDataClient,
		shopQueueClient,
	)
	cancelPurchaseClient := purchase.NewCancelPurchaseClient(
		sc_rdbcReadUserDataClient,
		smac_rdbcCancelPurchaseClient,
		shopQueueClient,
	)
	makePurchaseClient := purchase.NewMakePurchaseClient(
		makeShopAppraisalClient,
		sc_rdbcReadUserDataClient,
		smac_rdbcWritePurchaseClient,
		shopInventoryClient,
	)
	authClient := auth.NewAuthClient(
		authHashSetReaderClient,
		jwtClient,
		wc_mCharacterInfoClient,
	)

	// Proto Composition clients (Non-CFG)

	syncPBContractItemsClient := protoclient.
		NewPBContractItemsClient[*staticdb.SyncIndexMap](
		wc_ContractItemsClient,
	)
	localPBContractItemsClient := protoclient.
		NewPBContractItemsClient[*staticdb.LocalIndexMap](
		wc_ContractItemsClient,
	)
	syncPBGetBuybackAppraisalClient := protoclient.
		NewPBGetBuybackAppraisalClient[*staticdb.SyncIndexMap](
		wc_rdbcReadBuybackAppraisalClient,
	)
	localPBGetBuybackAppraisalClient := protoclient.
		NewPBGetBuybackAppraisalClient[*staticdb.LocalIndexMap](
		wc_rdbcReadBuybackAppraisalClient,
	)
	syncPBGetShopAppraisalClient := protoclient.
		NewPBGetShopAppraisalClient[*staticdb.SyncIndexMap](
		wc_rdbcReadShopAppraisalClient,
	)
	localPBGetShopAppraisalClient := protoclient.
		NewPBGetShopAppraisalClient[*staticdb.LocalIndexMap](
		wc_rdbcReadShopAppraisalClient,
	)
	syncPBNewBuybackAppraisalClient := protoclient.
		NewPBNewBuybackAppraisalClient[*staticdb.SyncIndexMap](
		makeBuybackAppraisalClient,
	)
	localPBNewBuybackAppraisalClient := protoclient.
		NewPBNewBuybackAppraisalClient[*staticdb.LocalIndexMap](
		makeBuybackAppraisalClient,
	)
	syncPBNewShopAppraisalClient := protoclient.
		NewPBNewShopAppraisalClient[*staticdb.SyncIndexMap](
		makeShopAppraisalClient,
	)
	localPBNewShopAppraisalClient := protoclient.
		NewPBNewShopAppraisalClient[*staticdb.LocalIndexMap](
		makeShopAppraisalClient,
	)
	statusBuybackAppraisalClient := protoclient.
		NewPBStatusBuybackAppraisalClient(
			localPBContractItemsClient,
			wc_ContractsClient,
			wc_StructureInfoClient,
		)
	statusShopAppraisalClient := protoclient.NewPBStatusShopAppraisalClient(
		localPBContractItemsClient,
		shopQueueClient,
		wc_StructureInfoClient,
	)
	pbBuybackContractQueueClient := protoclient.
		NewPBBuybackContractQueueClient(
			syncPBGetBuybackAppraisalClient,
			syncPBNewBuybackAppraisalClient,
			syncPBContractItemsClient,
			wc_ContractsClient,
			wc_StructureInfoClient,
		)
	pbShopContractQueueClient := protoclient.NewPBShopContractQueueClient(
		syncPBGetShopAppraisalClient,
		syncPBNewShopAppraisalClient,
		syncPBContractItemsClient,
		wc_ContractsClient,
		wc_StructureInfoClient,
	)
	pbShopPurchaseQueueClient := protoclient.NewPBShopPurchaseQueueClient(
		syncPBGetShopAppraisalClient,
		syncPBNewShopAppraisalClient,
		shopQueueClient,
	)
	pbShopInventoryClient := protoclient.NewPBShopInventoryClient(
		shopInventoryClient,
		shopPriceClient,
	)
	pbShopLocationsClient := protoclient.NewPBShopLocationsClient(
		wc_StructureInfoClient,
	)

	// Proto Composition clients (CFG)

	cfgGetBTypeMapsBuilderClient := protoclient.
		NewCfgGetBuybackSystemTypeMapsBuilderClient(
			webBTypeMapsBuilderReaderClient,
		)
	cfgMergeBTypeMapsBuilderClient := protoclient.
		NewCfgMergeBuybackSystemTypeMapsBuilderClient(
			webBTypeMapsBuilderReaderClient,
			webBTypeMapsBuilderWriterClient,
			webMarketsReaderClient,
		)
	cfgGetSTypeMapsBuilderClient := protoclient.
		NewCfgGetShopLocationTypeMapsBuilderClient(
			webSTypeMapsBuilderReaderClient,
		)
	cfgMergeSTypeMapsBuilderClient := protoclient.
		NewCfgMergeShopLocationTypeMapsBuilderClient(
			webSTypeMapsBuilderReaderClient,
			webSTypeMapsBuilderWriterClient,
			webMarketsReaderClient,
		)
	cfgGetBuybackSystemsClient := protoclient.NewCfgGetBuybackSystemsClient(
		webBuybackSystemsReaderClient,
	)
	cfgMergeBuybackSystemsClient := protoclient.
		NewCfgMergeBuybackSystemsClient(
			webBuybackSystemsReaderClient,
			webBuybackSystemsWriterClient,
			webBTypeMapsBuilderReaderClient,
			webSTypeMapsBuilderReaderClient,
		)
	cfgGetShopLocationsClient := protoclient.NewCfgGetShopLocationsClient(
		webShopLocationsReaderClient,
		wc_StructureInfoClient,
	)
	cfgMergeShopLocationsClient := protoclient.
		NewCfgMergeShopLocationsClient(
			webShopLocationsReaderClient,
			webShopLocationsWriterClient,
			webBTypeMapsBuilderReaderClient,
			webSTypeMapsBuilderReaderClient,
		)
	cfgGetMarketsClient := protoclient.NewCfgGetMarketsClient(
		webMarketsReaderClient,
		wc_StructureInfoClient,
	)
	cfgMergeMarketsClient := protoclient.NewCfgMergeMarketsClient(
		webMarketsReaderClient,
		webMarketsWriterClient,
		webBTypeMapsBuilderReaderClient,
		webSTypeMapsBuilderReaderClient,
	)

	return &Service{
		authClient:                     authClient,
		rReadAuthListClient:            authListReaderClient,
		rWriteAuthListClient:           authListWriterClient,
		rUserDataClient:                userDataClient,
		rdbcUserDataClient:             sc_rdbcReadUserDataClient,
		shopInventoryClient:            pbShopInventoryClient,
		shopContractQueueClient:        pbShopContractQueueClient,
		buybackContractQueueClient:     pbBuybackContractQueueClient,
		shopPurchaseQueueClient:        pbShopPurchaseQueueClient,
		shopMakePurchaseClient:         makePurchaseClient,
		shopCancelPurchaseClient:       cancelPurchaseClient,
		shopDeletePurchasesClient:      smac_rdbcDelPurchasesClient,
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
		shopLocationsClient:            pbShopLocationsClient,
	}
}
