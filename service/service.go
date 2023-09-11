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
	rCorpRawClient                 raw_.RawClient
	rMarketsRawClient              raw_.RawClient
	rStructureInfoRawClient        raw_.RawClient
	rAuthRawClient                 raw_.RawClient
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
	cfgGetBuybackBundleKeysClient  protoclient.CfgGetBuybackBundleKeysClient
	cfgGetShopBundleKeysClient     protoclient.CfgGetShopBundleKeysClient
	cfgGetMarketNamesClient        protoclient.CfgGetMarketNamesClient
	shopLocationsClient            protoclient.PBShopLocationsClient
	proto.UnimplementedEveTradingCoServer
}

func NewService(
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
	rBucketClient bucket.BucketClient,
	rRDBClient *rdb.RemoteDBClient,
	httpClient *http.Client,
) *Service {
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
		cCache,
		sCache,
	)
	jwtClient := jwt.NewJWTClient(
		unauthRawClient,
		authRawClient,
		cCache,
		sCache,
	)

	// Higher level bucket clients

	authHashSetReaderClient := bucketc.NewSC_AuthHashSetReaderClient(
		rBucketClient,
		sCache,
	)
	authHashSetWriterClient := bucketc.NewSAC_AuthHashSetWriterClient(
		rBucketClient,
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
			rBucketClient,
			sCache,
		)
	webBTypeMapsBuilderWriterClient :=
		bucketc.NewSMAC_WebBuybackSystemTypeMapsBuilderWriterClient(
			rBucketClient,
			webBTypeMapsBuilderReaderClient.GetAntiCache(),
		)
	webSTypeMapsBuilderReaderClient :=
		bucketc.NewSC_WebShopLocationTypeMapsBuilderReaderClient(
			rBucketClient,
			sCache,
		)
	webSTypeMapsBuilderWriterClient :=
		bucketc.NewSMAC_WebShopLocationTypeMapsBuilderWriterClient(
			rBucketClient,
			webSTypeMapsBuilderReaderClient.GetAntiCache(),
		)
	webBuybackSystemsReaderClient :=
		bucketc.NewSC_WebBuybackSystemsReaderClient(
			rBucketClient,
			sCache,
		)
	webBuybackSystemsWriterClient :=
		bucketc.NewSAC_WebBuybackSystemsWriterClient(
			rBucketClient,
			webBuybackSystemsReaderClient.GetAntiCache(),
		)
	webShopLocationsReaderClient :=
		bucketc.NewSC_WebShopLocationsReaderClient(
			rBucketClient,
			sCache,
		)
	webShopLocationsWriterClient :=
		bucketc.NewSAC_WebShopLocationsWriterClient(
			rBucketClient,
			webShopLocationsReaderClient.GetAntiCache(),
		)
	webMarketsReaderClient := bucketc.NewSC_WebMarketsReaderClient(
		rBucketClient,
		sCache,
	)
	webMarketsWriterClient := bucketc.NewSAC_WebMarketsWriterClient(
		rBucketClient,
		webMarketsReaderClient.GetAntiCache(),
	)
	webBuybackBundleKeysClient := bucketc.NewSC_WebBuybackBundleKeysClient(
		webBTypeMapsBuilderReaderClient,
		sCache,
	)
	webShopBundleKeysClient := bucketc.NewSC_WebShopBundleKeysClient(
		webSTypeMapsBuilderReaderClient,
		sCache,
	)

	// Higher Level remoteDB clients + Unreserved Location Assets

	wc_rdbcReadShopAppraisalClient := rdbc.NewWC_ReadShopAppraisalClient(
		rRDBClient,
		cCache,
		sCache,
	)
	wc_rdbcReadBuybackAppraisalClient :=
		rdbc.NewWC_ReadBuybackAppraisalClient(
			rRDBClient,
			cCache,
			sCache,
		)
	sc_rdbcReadShopQueueClient := rdbc.NewSC_ReadShopQueueClient(
		rRDBClient,
		sCache,
	)
	sc_rdbcReadUserDataClient := rdbc.NewSC_ReadUserDataClient(
		rRDBClient,
		sCache,
	)
	locationShopAssetsClient := locationassets.NewLocationShopAssetsClient(
		mAssetsCorporationClient,
		wc_rdbcReadShopAppraisalClient,
		cCache,
		sCache,
	)

	rdbcReadShopQueueAntiCache := sc_rdbcReadShopQueueClient.GetAntiCache()
	rdbcReadUserDataAntiCache := sc_rdbcReadUserDataClient.GetAntiCache()
	unreservedAntiCache := locationShopAssetsClient.GetUnreservedAntiCache()

	sac_rdbcWriteBuybackAppraisalClient :=
		rdbc.NewSAC_WriteBuybackAppraisalClient(
			rRDBClient,
			rdbcReadUserDataAntiCache,
		)
	smac_rdbcDelPurchasesClient := rdbc.NewSMAC_DelPurchasesClient(
		rRDBClient,
		rdbcReadShopQueueAntiCache,
		unreservedAntiCache,
	)
	smac_rdbcCancelPurchaseClient := rdbc.NewSMAC_CancelPurchaseClient(
		rRDBClient,
		rdbcReadUserDataAntiCache,
		rdbcReadShopQueueAntiCache,
		unreservedAntiCache,
	)
	smac_rdbcWritePurchaseClient := rdbc.NewSMAC_WritePurchaseClient(
		rRDBClient,
		rdbcReadUserDataAntiCache,
		rdbcReadShopQueueAntiCache,
		unreservedAntiCache,
	)

	// Non-Proto Composition clients

	marketPriceClient := marketprice.NewMarketPriceClient(
		mOrdersRegionClient,
		mOrdersStructureClient,
		cCache,
		sCache,
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
		cCache,
		sCache,
	)
	wc_ContractsClient := contracts.NewWC_ContractsClient(
		mContractsCorporationClient,
		cCache,
		sCache,
	)
	wc_ContractItemsClient := contracts.NewWC_SingleContractItemsClient(
		mContractItemsClient,
		cCache,
		sCache,
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
	cfgGetBuybackBundleKeysClient := protoclient.NewCfgGetBuybackBundleKeysClient(
		webBuybackBundleKeysClient,
	)
	cfgGetShopBundleKeysClient := protoclient.NewCfgGetShopBundleKeysClient(
		webShopBundleKeysClient,
	)
	cfgGetMarketNamesClient := protoclient.NewCfgGetMarketNamesClient(
		webMarketsReaderClient,
	)

	return &Service{
		rCorpRawClient:                 corpRawClient,
		rMarketsRawClient:              marketsRawClient,
		rStructureInfoRawClient:        structureInfoRawClient,
		rAuthRawClient:                 authRawClient,
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
		cfgGetBuybackBundleKeysClient:  cfgGetBuybackBundleKeysClient,
		cfgGetShopBundleKeysClient:     cfgGetShopBundleKeysClient,
		cfgGetMarketNamesClient:        cfgGetMarketNamesClient,
		shopLocationsClient:            pbShopLocationsClient,
	}
}
