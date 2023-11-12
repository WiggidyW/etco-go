package esi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
)

const (
	CHARACTER_INFO_BUF_CAP          int           = 0
	CHARACTER_INFO_LOCK_TTL         time.Duration = 30 * time.Second
	CHARACTER_INFO_LOCK_MAX_BACKOFF time.Duration = 10 * time.Second
	CHARACTER_INFO_MIN_EXPIRES_IN   time.Duration = 24 * time.Hour
	CHARACTER_INFO_METHOD           string        = http.MethodGet
)

func init() {
	keys.TypeStrCharacterInfo = cache.RegisterType[CharacterInfo]("jwks", CHARACTER_INFO_BUF_CAP)
}

type CharacterInfo struct {
	AllianceId    *int32 `json:"alliance_id,omitempty"`
	CorporationId int32  `json:"corporation_id"`
	Name          string `json:"name"`
	// Birthday       time.Time `json:"birthday"`
	// BloodlineId    int32     `json:"bloodline_id"`
	// Description    *string   `json:"description"`
	// FactionId      *int32    `json:"faction_id"`
	// Gender         string    `json:"gender"`
	// RaceId         int32     `json:"race_id"`
	// SecurityStatus *float64  `json:"security_status"`
	// Title          *string   `json:"title"`
}

func characterInfoUrl(characterId int32) string {
	return fmt.Sprintf(
		"%s/characters/%d/?datasource=%s",
		BASE_URL,
		characterId,
		DATASOURCE,
	)
}

func GetCharacterInfo(x cache.Context, characterId int32) (
	rep *CharacterInfo,
	expires time.Time,
	err error,
) {
	return infoGet(
		x,
		characterInfoUrl(characterId),
		CHARACTER_INFO_METHOD,
		keys.TypeStrCharacterInfo,
		keys.CacheKeyCharacterInfo(characterId),
		CHARACTER_INFO_LOCK_TTL,
		CHARACTER_INFO_LOCK_MAX_BACKOFF,
		CHARACTER_INFO_MIN_EXPIRES_IN,
		nil,
		entityInfoHandleErr[CharacterInfo],
	)
}
