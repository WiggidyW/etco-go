package keys

var (
	// namespace (multi-cache set, idk how to explain tbh)
	TypeStrNamespace string

	// esi
	TypeStrJWKS               string
	TypeStrCharacterInfo      string
	TypeStrCorporationInfo    string
	TypeStrAllianceInfo       string
	TypeStrStructureInfo      string
	TypeStrAuthToken          string
	TypeStrCorpToken          string
	TypeStrStructureInfoToken string
	TypeStrMarketsToken       string

	// bucket
	TypeStrWebBuybackSystemTypeMapsBuilder string
	TypeStrWebBuybackBundleKeys            string
	TypeStrWebShopLocationTypeMapsBuilder  string
	TypeStrWebShopBundleKeys               string
	TypeStrWebBuybackSystems               string
	TypeStrWebShopLocations                string
	TypeStrWebMarkets                      string
	TypeStrBuildConstData                  string
	TypeStrAuthHashSet                     string

	// remoteDB
	TypeStrNSUserData                string
	TypeStrUserBuybackAppraisalCodes string
	TypeStrUserShopAppraisalCodes    string
	TypeStrUserCancelledPurchase     string
	TypeStrUserMadePurchase          string

	TypeStrBuybackAppraisal string
	TypeStrShopAppraisal    string

	TypeStrRawPurchaseQueue string

	// // composition
	TypeStrPurchaseQueue         string
	TypeStrLocationPurchaseQueue string

	TypeStrContractItems string

	TypeStrContracts string

	TypeStrNSRegionMarketOrders    string
	TypeStrRegionMarketOrders      string
	TypeStrNSStructureMarketOrders string
	TypeStrStructureMarketOrders   string

	TypeStrNSRawShopAssets      string
	TypeStrRawShopAssets        string
	TypeStrUnreservedShopAssets string

	TypeStrTokenCharacter string
)
