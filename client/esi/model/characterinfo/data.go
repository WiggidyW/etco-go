package characterinfo

type CharacterInfoModel struct {
	AllianceId *int32 `json:"alliance_id,omitempty"`
	// Birthday       time.Time `json:"birthday"`
	// BloodlineId    int32     `json:"bloodline_id"`
	CorporationId int32 `json:"corporation_id"`
	// Description    *string   `json:"description,omitempty"`
	// FactionId      *int32    `json:"faction_id,omitempty"`
	// Gender         string    `json:"gender"`
	Name string `json:"name"`
	// RaceId         int32     `json:"race_id"`
	// SecurityStatus *float64  `json:"security_status,omitempty"`
	// Title          *string   `json:"title,omitempty"`
}
