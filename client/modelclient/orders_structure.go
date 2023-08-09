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
	ORDERS_STRUCTURE_USE_AUTH         bool  = true
	ORDERS_STRUCTURE_ENTRIES_PER_PAGE int32 = 1000
)

type ClientOrdersStructure struct {
	naiveclient.NaivePagesClient[
		EntryOrdersStructure,
		UrlParamsOrdersStructure,
	]
}

func NewClientOrdersStructure(
	rawClient *rawclient.RawClient,
	minExpires time.Duration,
	modelServerLockTTL time.Duration,
	modelServerLockMaxWait time.Duration,
	headPool *cache.BufferPool,
	headServerLockTTL time.Duration,
	headServerLockMaxWait time.Duration,
	clientCache cache.SharedClientCache,
	serverCache cache.SharedServerCache,
) *ClientOrdersStructure {
	return &ClientOrdersStructure{
		naiveclient.NewNaivePagesClient[
			EntryOrdersStructure,
			UrlParamsOrdersStructure,
		](
			rawClient,
			ASSETS_CORPORATION_USE_AUTH,
			ASSETS_CORPORATION_ENTRIES_PER_PAGE,
			minExpires,
			cache.NewBufferPool(0),
			modelServerLockTTL,
			modelServerLockMaxWait,
			headPool,
			headServerLockTTL,
			headServerLockMaxWait,
			clientCache,
			serverCache,
		),
	}
}

func NewFetchParamsOrdersStructure(
	structureId int64,
	refreshToken string,
) *naiveclient.NaiveClientFetchParams[UrlParamsOrdersStructure] {
	return naiveclient.NewNaiveClientFetchParams[UrlParamsOrdersStructure](
		UrlParamsOrdersStructure{
			structureId: structureId,
		},
		&refreshToken,
		nil,
	)
}

// type ModelOrdersStructure = []EntryItemOrdersStructure

type EntryOrdersStructure struct {
	// Duration   int32 `json:"duration"`
	IsBuyOrder bool `json:"is_buy_order"`
	// Issued       time.Time `json:"issued"`
	// LocationId int64 `json:"location_id"`
	// MinVolume   int32 `json:"min_volume"`
	// OrderId     int64 `json:"order_id"`
	Price float64 `json:"price"`
	// Range       string `json:"range"`
	TypeId       int32 `json:"type_id"`
	VolumeRemain int32 `json:"volume_remain"`
	// VolumeTotal int32 `json:"volume_total"`
}

type UrlParamsOrdersStructure struct {
	structureId int64
}

func (p UrlParamsOrdersStructure) PageCacheKey(page *int32) string {
	query := fmt.Sprintf(
		"%d/?datasource=%s",
		p.structureId,
		DATASOURCE,
	)
	query = addQueryInt32(query, "page", page)
	return query
}
func (p UrlParamsOrdersStructure) PageUrl(page *int32) string {
	return fmt.Sprintf("%s/markets/structures/%s", BASE_URL, p.PageCacheKey(page))
}
func (p UrlParamsOrdersStructure) CacheKey() string {
	return p.PageCacheKey(nil)
}
func (p UrlParamsOrdersStructure) Url() string {
	return p.PageUrl(nil)
}
func (UrlParamsOrdersStructure) Method() string {
	return http.MethodGet
}
