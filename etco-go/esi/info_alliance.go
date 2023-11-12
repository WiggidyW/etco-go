package esi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
)

const (
	ALLIANCE_INFO_BUF_CAP          int           = 0
	ALLIANCE_INFO_LOCK_TTL         time.Duration = 30 * time.Second
	ALLIANCE_INFO_LOCK_MAX_BACKOFF time.Duration = 10 * time.Second
	ALLIANCE_INFO_MIN_EXPIRES_IN   time.Duration = 24 * time.Hour
	ALLIANCE_INFO_METHOD           string        = http.MethodGet
)

func init() {
	keys.TypeStrAllianceInfo = cache.RegisterType[AllianceInfo]("jwks", ALLIANCE_INFO_BUF_CAP)
}

type AllianceInfo struct {
	Name   string `json:"name"`
	Ticker string `json:"ticker"`
	// CreatorCorporationId int32 `json:"creator_corporation_id"`
	// CreatorId            int32 `json:"creator_id"`
	// DateFounded          time.Time `json:"date_founded"`
	// ExecutorCorporationId *int32 `json:"executor_corporation_id"`
	// FactionId            *int32 `json:"faction_id"`
}

func allianceInfoUrl(allianceId int32) string {
	return fmt.Sprintf(
		"%s/alliances/%d/?datasource=%s",
		BASE_URL,
		allianceId,
		DATASOURCE,
	)
}

func GetAllianceInfo(x cache.Context, allianceId int32) (
	rep *AllianceInfo,
	expires time.Time,
	err error,
) {
	return infoGet(
		x,
		allianceInfoUrl(allianceId),
		ALLIANCE_INFO_METHOD,
		keys.TypeStrAllianceInfo,
		keys.CacheKeyAllianceInfo(allianceId),
		ALLIANCE_INFO_LOCK_TTL,
		ALLIANCE_INFO_LOCK_MAX_BACKOFF,
		ALLIANCE_INFO_MIN_EXPIRES_IN,
		nil,
		entityInfoHandleErr[AllianceInfo],
	)
}
