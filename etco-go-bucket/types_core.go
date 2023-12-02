package etcogobucket

type CoreBucketData struct {
	BuybackSystemTypeMaps []BuybackSystemTypeMap
	ShopLocationTypeMaps  []ShopLocationTypeMap
	BuybackSystems        map[SystemId]BuybackSystem
	ShopLocations         map[LocationId]ShopLocation
	BannedFlagSets        []BannedFlagSet
	Pricings              []Pricing
	Markets               []Market
	HaulRouteTypeMaps     []HaulRouteTypeMap
	HaulRoutes            map[HaulRouteSystemsKey]HaulRouteInfoIndex
	HaulRouteInfos        []HaulRouteInfo
	UpdaterData           CoreUpdaterData
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

// type BuybackSystemTypeMaps = []BuybackSystemTypeMap
type BuybackSystemTypeMap = map[TypeId]BuybackTypePricing
type BuybackTypePricing struct {
	ReprocessingEfficiency uint8 // 0 = nil, 1 - 100 = efficiency
	PricingIndex           int   // nil -> -1
}

// type ShopLocationTypeMaps = []ShopLocationTypeMap
type ShopLocationTypeMap = map[TypeId]ShopTypePricing
type ShopTypePricing = int // PricingIndex

// type HaulRouteTypeMaps = []HaulRouteTypeMap
type HaulRouteTypeMap = map[TypeId]HaulRouteTypePricing
type HaulRouteTypePricing = int // PricingIndex

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

// type HaulRoutes = map[[4]byte]HaulRouteInfoIndex
type HaulRouteSystemsKey = [4]byte // [StartSystemIndex, EndSystemIndex]
type HaulRouteInfoIndex = uint16

// type HaulRouteInfos = []HaulRouteInfo
type HaulRouteInfo struct {
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

	TypeMapIndex uint8
}

func (hri HaulRouteInfo) MaxVolume() Scientific16 {
	return Scientific16{
		Base:   hri.MaxVolumeS16Base,
		Zeroes: hri.MaxVolumeS16Zeroes,
	}
}
func (hri HaulRouteInfo) MinReward() Scientific16 {
	return Scientific16{
		Base:   hri.MinRewardS16Base,
		Zeroes: hri.MinRewardS16Zeroes,
	}
}
func (hri HaulRouteInfo) MaxReward() Scientific16 {
	return Scientific16{
		Base:   hri.MaxRewardS16Base,
		Zeroes: hri.MaxRewardS16Zeroes,
	}
}

type HaulRouteRewardStrategy uint8

const (
	// lesser of (m3Fee * volume) and (collateral * collateralRate)
	HRRSLesserOf HaulRouteRewardStrategy = iota
	// greater of (m3Fee * volume) and (collateral * collateralRate)
	HRRSGreaterOf
	// (m3Fee * volume) + (collateral * collateralRate)
	HRRSSum
)
