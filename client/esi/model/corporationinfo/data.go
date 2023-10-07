package corporationinfo

type CorporationInfoModel struct {
	AllianceId *int32 `json:"alliance_id,omitempty"`
	// CeoId      int32 `json:"ceo_id"`
	// CreatorId  int32 `json:"creator_id"`
	// DateFounded *time.Time `json:"date_founded"`
	// Description *string `json:"description,omitempty"`
	// FactionId   *int32 `json:"faction_id,omitempty"`
	// HomeStationId *int32 `json:"home_station_id,omitempty"`
	// MemberCount int32 `json:"member_count"`
	Name string `json:"name"`
	// Shares      int64 `json:"shares"`
	// TaxRate     float64 `json:"tax_rate"`
	Ticker string `json:"ticker"`
	// Url         *string `json:"url"`
	// WarEligible *bool `json:"war_eligible"`
}
