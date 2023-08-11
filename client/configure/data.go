package configure

type ShopLocationMap = map[int64]ShopLocation

type BuybackSystemMap = map[int32]BuybackSystem

type ShopLocation struct {
	TypeMap     string   `json:"typeMap"` // type map key
	BannedFlags []string `json:"bannedFlags"`
}

type BuybackSystem struct {
	TypeMap string   `json:"typeMap"` // type map key
	M3Fee   *float64 `json:"m3Fee"`
}

type BuybackTypeMap = map[int32]BuybackType

type ShopTypeMap = map[int32]ShopType

type TypeKey = string

type BuybackType = map[TypeKey]*BuybackTypeValue

type ShopType = map[TypeKey]*ShopTypeValue

type BuybackTypeValue struct {
	ReprocessingEfficiency uint8    `json:"reprocessingEfficiency"`
	Pricing                *Pricing `json:"pricing"`
}

type ShopTypeValue = Pricing

type Pricing struct {
	IsBuy      bool   `json:"isBuy"`
	Percentile uint8  `json:"percentile"`
	Modifier   uint8  `json:"modifier"`
	Market     string `json:"market"` // market key
}

type MarketMap = map[string]Market

type Market struct {
	RefreshToken *string `json:"refreshToken"`
	LocationId   int64   `json:"locationId"`
	IsStructure  bool    `json:"isStructure"`
}

func MergeBuybackTypeMap(
	old BuybackTypeMap,
	changes BuybackTypeMap,
) BuybackTypeMap {
	if changes == nil {
		return old
	}

	for system, type_ := range changes {
		for typeKey, typeValue := range type_ {
			if _, ok := old[system]; !ok {
				old[system] = make(BuybackType)
			}
			old[system][typeKey] = typeValue
		}
	}
	return old
}
