package modelclient

import (
	"fmt"
	"net/http"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client/naiveclient"
	"github.com/WiggidyW/weve-esi/client/naiveclient/rawclient"
)

const (
	CHARACTER_INFO_USE_AUTH bool = false
)

type ClientCharacterInfo struct {
	naiveclient.CachingNaivePageModelClient[
		ModelCharacterInfo,
		UrlParamsCharacterInfo,
	]
}

func NewClientCharacterInfo(
	rawClient *rawclient.RawClient,
	minExpires time.Duration,
	modelServerLockTTL time.Duration,
	modelServerLockMaxWait time.Duration,
	clientCache cache.SharedClientCache,
	serverCache cache.SharedServerCache,
) *ClientCharacterInfo {
	cachingClient := naiveclient.NewCachingNaivePageModelClient[
		ModelCharacterInfo,
		UrlParamsCharacterInfo,
	](
		rawClient,
		STRUCTURE_INFO_USE_AUTH,
		minExpires,
		cache.NewBufferPool(0),
		clientCache,
		serverCache,
		modelServerLockTTL,
		modelServerLockMaxWait,
	)
	return &ClientCharacterInfo{cachingClient}
}

func NewFetchParamsCharacterInfo(
	characterId int32,
) *naiveclient.NaiveClientFetchParams[UrlParamsCharacterInfo] {
	return naiveclient.NewNaiveClientFetchParams[UrlParamsCharacterInfo](
		UrlParamsCharacterInfo{characterId},
		nil,
		nil,
	)
}

type ModelCharacterInfo struct {
	AllianceId *int32 `json:"alliance_id,omitempty"`
	// Birthday       time.Time `json:"birthday"`
	// BloodlineId    int32     `json:"bloodline_id"`
	CorporationId int32 `json:"corporation_id"`
	// Description    *string   `json:"description,omitempty"`
	// FactionId      *int32    `json:"faction_id,omitempty"`
	// Gender         string    `json:"gender"`
	// Name           string    `json:"name"`
	// RaceId         int32     `json:"race_id"`
	// SecurityStatus *float64  `json:"security_status,omitempty"`
	// Title          *string   `json:"title,omitempty"`
}

type UrlParamsCharacterInfo struct {
	characterId int32
}

func (p UrlParamsCharacterInfo) CacheKey() string {
	return fmt.Sprintf(
		"%d/?datasource=%s",
		p.characterId,
		DATASOURCE,
	)
}

func (p UrlParamsCharacterInfo) Url() string {
	return fmt.Sprintf("%s/characters/%s", BASE_URL, p.CacheKey())
}

func (UrlParamsCharacterInfo) Method() string {
	return http.MethodGet
}
