package main

import (
	"fmt"
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
	"github.com/WiggidyW/etco-go/client/purchase"
	rdbc "github.com/WiggidyW/etco-go/client/remotedb"
	"github.com/WiggidyW/etco-go/client/shopqueue"
	"github.com/WiggidyW/etco-go/client/structureinfo"
	"github.com/WiggidyW/etco-go/client/userdata"
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

func main() {
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
	WC_StructureInfoClient := structureinfo.NewWC_StructureInfoClient(
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

	// Proto Composition clients

	// flag.Parse()
	// env.Init()

	// if *runLocal {
	// 	fmt.Println("Starting server in local mode")
	// } else {
	// 	fmt.Println("Starting server")

	fmt.Println("Hello, code inspector!")

	// gob.Register(map[string]interface{}{})

	// service := client.NewClient(*runLocal)
	// server := pb.NewWeveEsiServer(service)

	// corsWrapper := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// 	AllowedMethods: []string{"POST"},
	// 	AllowedHeaders: []string{"Content-Type"},
	// })
	// handler := corsWrapper.Handler(server)

	// http.ListenAndServe(env.LISTEN_ADDRESS, handler)
}
