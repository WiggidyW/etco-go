package configure

const (
	B_TYPE_MAPS_BUILDER_DOMAIN_KEY string = "btypemapsbuilder"
	S_TYPE_MAPS_BUILDER_DOMAIN_KEY string = "stypemapsbuilder"
	MARKETS_DOMAIN_KEY             string = "markets"
	SHOP_LOCATIONS_DOMAIN_KEY      string = "shoplocations"
	BUYBACK_SYSTEMS_DOMAIN_KEY     string = "buybacksystems"
)

type BundleKey = string
type MarketName = string
type LocationId = int64
type SystemId = int32
type TypeId = int32

type Markets = map[MarketName]Market
type Market struct {
	RefreshToken *string
	LocationId   int64
	IsStructure  bool
}

type ShopLocations = map[LocationId]ShopLocation
type ShopLocation struct {
	BundleKey   string
	BannedFlags []string
}

type BuybackSystems = map[SystemId]BuybackSystem
type BuybackSystem struct {
	BundleKey BundleKey
	M3Fee     float64
}

type TypePricing struct {
	IsBuy      bool
	Percentile uint8
	Modifier   uint8
	Market     MarketName // keys to markets
}

type ShopTypePricing = TypePricing
type ShopLocationTypeBundle = map[BundleKey]ShopTypePricing
type ShopLocationTypeMapsBuilder = map[TypeId]ShopLocationTypeBundle

type BuybackTypePricing struct {
	Pricing *TypePricing
	ReprEff uint8
}
type BuybackSystemTypeBundle = map[BundleKey]BuybackTypePricing
type BuybackSystemTypeMapsBuilder = map[TypeId]BuybackSystemTypeBundle

type PBMergeResponse struct {
	Modified   bool
	MergeError error
}
