package assetscorporation

type AssetsCorporationEntry struct {
	// IsBlueprintCopy *bool  `json:"is_blueprint_copy,omitempty"`
	// IsSingleton     bool   `json:"is_singleton"`
	ItemId       int64  `json:"item_id"`
	LocationFlag string `json:"location_flag"`
	LocationId   int64  `json:"location_id"`
	// LocationType    string `json:"location_type"`
	Quantity int32 `json:"quantity"`
	TypeId   int32 `json:"type_id"`
}
