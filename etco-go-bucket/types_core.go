package etcogobucket

type CoreBucketData struct {
	BuybackSystemTypeMaps []BuybackSystemTypeMap
	ShopLocationTypeMaps  []ShopLocationTypeMap
	BuybackSystems        map[SystemId]BuybackSystem
	ShopLocations         map[LocationId]ShopLocation
	BannedFlagSets        []BannedFlagSet
	Pricings              []Pricing
	Markets               []Market
}

// type BuybackSystemTypeMaps = []BuybackSystemTypeMap
type BuybackSystemTypeMap = map[TypeId]BuybackTypePricing
type BuybackTypePricing struct {
	ReprocessingEfficiency uint8 // 0 = nil, 1 - 100 = efficiency
	PricingIndex           int   // nil -> -1
}

// type ShopLocationTypeMaps = []ShopLocationTypeMap
type ShopLocationTypeMap = map[TypeId]ShopTypePricing
type ShopTypePricing = int // PricingIndex

// type BuybackSystems = map[SystemId]BuybackSystem
type BuybackSystem struct {
	M3Fee        float64
	TaxRate      float64 // 0-1
	TypeMapIndex int
}

// type ShopLocations = map[LocationId]ShopLocation
type ShopLocation struct {
	BannedFlagSetIndex int     // nil -> -1
	TaxRate            float64 // 0-1
	TypeMapIndex       int
}

// type BannedFlagSets = []BannedFlagSet
type BannedFlagSet = map[BannedFlag]struct{} // set of banned flags
type BannedFlag = string                     // english

// type Markets = []Market
type Market struct {
	Name         string // user-defined
	RefreshToken *string
	LocationId   int64
	IsStructure  bool
}

// type Pricings = []Pricing
type Pricing struct {
	IsBuy       bool
	Percentile  uint8 // 0 - 100
	Modifier    uint8 // 1 - 255
	MarketIndex int
}
