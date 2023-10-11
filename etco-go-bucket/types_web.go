package etcogobucket

type WebBucketData struct {
	BuybackSystemTypeMapsBuilder map[TypeId]WebBuybackSystemTypeBundle
	ShopLocationTypeMapsBuilder  map[TypeId]WebShopLocationTypeBundle
	BuybackSystems               map[SystemId]WebBuybackSystem
	ShopLocations                map[LocationId]WebShopLocation
	Markets                      map[MarketName]WebMarket
}

// type WebBuybackSystemTypeMapsBuilder = map[TypeId]WebBuybackSystemTypeBundle
type WebBuybackSystemTypeBundle = map[BundleKey]WebBuybackTypePricing
type WebBuybackTypePricing struct {
	Pricing                *WebTypePricing
	ReprocessingEfficiency uint8
}

// type WebShopLocationTypeMapsBuilder = map[TypeId]WebShopLocationTypeBundle
type WebShopLocationTypeBundle = map[BundleKey]WebShopTypePricing
type WebShopTypePricing = WebTypePricing

// type WebBuybackSystems = map[SystemId]WebBuybackSystem
type WebBuybackSystem struct {
	BundleKey BundleKey
	M3Fee     float64
}

// type WebShopLocations = map[LocationId]WebShopLocation
type WebShopLocation struct {
	BundleKey   string
	BannedFlags []string
}

// type WebMarkets = map[MarketName]WebMarket
type WebMarket struct {
	RefreshToken *string
	LocationId   int64
	IsStructure  bool
}
type MarketName = string

// shared
type WebTypePricing struct {
	IsBuy      bool
	Percentile uint8
	Modifier   uint8
	MarketName MarketName // keys to markets
}
type BundleKey = string
