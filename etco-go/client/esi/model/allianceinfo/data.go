package allianceinfo

type AllianceInfoModel struct {
	// CreatorCorporationId int32 `json:"creator_corporation_id"`
	// CreatorId            int32 `json:"creator_id"`
	// DateFounded          time.Time `json:"date_founded"`
	// ExecutorCorporationId *int32 `json:"executor_corporation_id"`
	// FactionId            *int32 `json:"faction_id"`
	Name   string `json:"name"`
	Ticker string `json:"ticker"`
}
