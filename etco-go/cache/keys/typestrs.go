package keys

var (
	// esi
	TypeStrJWKS               Key
	TypeStrCharacterInfo      Key
	TypeStrCorporationInfo    Key
	TypeStrAllianceInfo       Key
	TypeStrStructureInfo      Key
	TypeStrAuthToken          Key
	TypeStrCorpToken          Key
	TypeStrStructureInfoToken Key
	TypeStrMarketsToken       Key

	// bucket
	TypeStrWebBuybackSystemTypeMapsBuilder Key
	TypeStrWebBuybackBundleKeys            Key
	TypeStrWebShopLocationTypeMapsBuilder  Key
	TypeStrWebShopBundleKeys               Key
	TypeStrWebBuybackSystems               Key
	TypeStrWebShopLocations                Key
	TypeStrWebMarkets                      Key
	TypeStrBuildConstData                  Key
	TypeStrAuthHashSet                     Key

	// remoteDB
	TypeStrNSUserData                Key
	TypeStrUserBuybackAppraisalCodes Key
	TypeStrUserShopAppraisalCodes    Key
	TypeStrUserCancelledPurchase     Key
	TypeStrUserMadePurchase          Key
	TypeStrPrevContracts             Key

	TypeStrBuybackAppraisal Key
	TypeStrShopAppraisal    Key

	TypeStrRawPurchaseQueue Key

	// // composition
	TypeStrPurchaseQueue         Key
	TypeStrLocationPurchaseQueue Key

	TypeStrContractItems Key

	TypeStrNSContracts      Key
	TypeStrBuybackContracts Key
	TypeStrShopContracts    Key

	TypeStrNSRegionMarketOrders    Key
	TypeStrRegionMarketOrders      Key
	TypeStrNSStructureMarketOrders Key
	TypeStrStructureMarketOrders   Key

	TypeStrNSRawShopAssets      Key
	TypeStrRawShopAssets        Key
	TypeStrUnreservedShopAssets Key

	TypeStrTokenCharacter Key
)
