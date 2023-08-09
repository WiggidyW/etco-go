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
	STRUCTURE_INFO_USE_AUTH bool = true
)

type ClientStructureInfo struct {
	naiveclient.CachingNaivePageModelClient[
		ModelStructureInfo,
		UrlParamsStructureInfo,
	]
}

func NewClientStructureInfo(
	rawClient *rawclient.RawClient,
	minExpires time.Duration,
	modelServerLockTTL time.Duration,
	modelServerLockMaxWait time.Duration,
	clientCache cache.SharedClientCache,
	serverCache cache.SharedServerCache,
) *ClientStructureInfo {
	cachingClient := naiveclient.NewCachingNaivePageModelClient[
		ModelStructureInfo,
		UrlParamsStructureInfo,
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
	return &ClientStructureInfo{cachingClient}
}

func NewFetchParamsStructureInfo(
	structureId int64,
	refreshToken string,
) *naiveclient.NaiveClientFetchParams[UrlParamsStructureInfo] {
	return naiveclient.NewNaiveClientFetchParams[UrlParamsStructureInfo](
		UrlParamsStructureInfo{
			structureId: structureId,
		},
		&refreshToken,
		nil,
	)
}

type ModelStructureInfo struct {
	Name          string            `json:"name"`
	OwnerId       int32             `json:"owner_id"`
	Position      StructurePosition `json:"position"`
	SolarSystemId int32             `json:"solar_system_id"`
	TypeId        *int32            `json:"type_id,omitempty"`
}

type StructurePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type UrlParamsStructureInfo struct {
	structureId int64
}

func (p UrlParamsStructureInfo) CacheKey() string {
	return fmt.Sprintf(
		"%d/?datasource=%s",
		p.structureId,
		DATASOURCE,
	)
}

func (p UrlParamsStructureInfo) Url() string {
	return fmt.Sprintf("%s/universe/structures/%s", BASE_URL, p.CacheKey())
}

func (UrlParamsStructureInfo) Method() string {
	return http.MethodGet
}
