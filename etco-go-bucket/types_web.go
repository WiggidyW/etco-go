package etcogobucket

type WebBucketData struct {
	BuybackSystemTypeMapsBuilder map[TypeId]WebBuybackSystemTypeBundle
	ShopLocationTypeMapsBuilder  map[TypeId]WebShopLocationTypeBundle
	HaulRouteTypeMapsBuilder     map[TypeId]WebHaulRouteTypeBundle
	BuybackSystems               map[SystemId]WebBuybackSystem
	ShopLocations                map[LocationId]WebShopLocation
	HaulRoutes                   map[WebHaulRouteSystemsKey]WebHaulRoute
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
	TaxRate   float64 // 0-1
	M3Fee     float64
}

// type WebShopLocations = map[LocationId]WebShopLocation
type WebShopLocation struct {
	BundleKey   string
	TaxRate     float64 // 0-1
	BannedFlags []string
}

// type WebMarkets = map[MarketName]WebMarket
type WebMarket struct {
	RefreshToken *string
	LocationId   int64
	IsStructure  bool
}
type MarketName = string

// type WebHaulRouteTypeMapsBuilder = map[TypeId]WebHaulRouteTypeBundle
type WebHaulRouteTypeBundle = map[BundleKey]WebHaulRouteTypePricing
type WebHaulRouteTypePricing = WebTypePricing

type WebHaulRouteSystemsKey = [8]byte // [StartSystemId, EndSystemId]
type WebHaulRoute struct {
	// // calculates reward
	// Inner structs make Gobs bigger, so we store Scientific16s inline

	// MaxVolume 	  Scientific16
	MaxVolumeS16Base   uint8
	MaxVolumeS16Zeroes uint8

	// MinReward      Scientific16
	MinRewardS16Base   uint8
	MinRewardS16Zeroes uint8

	// MaxReward 	   Scientific16
	MaxRewardS16Base   uint8
	MaxRewardS16Zeroes uint8

	TaxRate DecPercentage

	M3Fee          uint16
	CollateralRate DecPercentage
	RewardStrategy HaulRouteRewardStrategy

	// // calculates collateral

	BundleKey BundleKey
}

func (whr WebHaulRoute) MaxVolume() Scientific16 {
	return Scientific16{
		Base:   whr.MaxVolumeS16Base,
		Zeroes: whr.MaxVolumeS16Zeroes,
	}
}
func (whr WebHaulRoute) MinReward() Scientific16 {
	return Scientific16{
		Base:   whr.MinRewardS16Base,
		Zeroes: whr.MinRewardS16Zeroes,
	}
}
func (whr WebHaulRoute) MaxReward() Scientific16 {
	return Scientific16{
		Base:   whr.MaxRewardS16Base,
		Zeroes: whr.MaxRewardS16Zeroes,
	}
}

// shared
type WebTypePricing struct {
	IsBuy      bool
	Percentile uint8
	Modifier   uint8
	MarketName MarketName // keys to markets
}
type BundleKey = string
