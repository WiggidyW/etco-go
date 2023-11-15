package esi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
)

const (
	CORPORATION_INFO_BUF_CAP          int           = 0
	CORPORATION_INFO_LOCK_TTL         time.Duration = 30 * time.Second
	CORPORATION_INFO_LOCK_MAX_BACKOFF time.Duration = 10 * time.Second
	CORPORATION_INFO_MIN_EXPIRES_IN   time.Duration = 24 * time.Hour
	CORPORATION_INFO_METHOD           string        = http.MethodGet
)

func init() {
	keys.TypeStrCorporationInfo = cache.RegisterType[CorporationInfo]("jwks", CORPORATION_INFO_BUF_CAP)
}

type CorporationInfo struct {
	AllianceId *int32 `json:"alliance_id"`
	Name       string `json:"name"`
	Ticker     string `json:"ticker"`
	// CeoId      int32 `json:"ceo_id"`
	// CreatorId  int32 `json:"creator_id"`
	// DateFounded *time.Time `json:"date_founded"`
	// Description *string `json:"description"`
	// FactionId   *int32 `json:"faction_id"`
	// HomeStationId *int32 `json:"home_station_id"`
	// MemberCount int32 `json:"member_count"`
	// Shares      int64 `json:"shares"`
	// TaxRate     float64 `json:"tax_rate"`
	// Url         *string `json:"url"`
	// WarEligible *bool `json:"war_eligible"`
}

func corporationInfoUrl(corporationId int32) string {
	return fmt.Sprintf(
		"%s/corporations/%d/?datasource=%s",
		BASE_URL,
		corporationId,
		DATASOURCE,
	)
}

func GetCorporationInfo(x cache.Context, corporationId int32) (
	rep *CorporationInfo,
	expires time.Time,
	err error,
) {
	return infoGet(
		x,
		corporationInfoUrl(corporationId),
		CORPORATION_INFO_METHOD,
		keys.TypeStrCorporationInfo,
		keys.CacheKeyCorporationInfo(corporationId),
		CORPORATION_INFO_LOCK_TTL,
		CORPORATION_INFO_LOCK_MAX_BACKOFF,
		CORPORATION_INFO_MIN_EXPIRES_IN,
		nil,
		entityInfoHandleErr[CorporationInfo],
	)
}
