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
	ASSETS_CORPORATION_USE_AUTH         bool  = true
	ASSETS_CORPORATION_ENTRIES_PER_PAGE int32 = 1000
)

type ClientAssetsCorporation struct {
	naiveclient.NaivePagesClient[
		EntryAssetsCorporation,
		UrlParamsAssetsCorporation,
	]
}

func NewClientAssetsCorporation(
	rawClient *rawclient.RawClient,
	minExpires time.Duration,
	modelServerLockTTL time.Duration,
	modelServerLockMaxWait time.Duration,
	headPool *cache.BufferPool,
	headServerLockTTL time.Duration,
	headServerLockMaxWait time.Duration,
	clientCache cache.SharedClientCache,
	serverCache cache.SharedServerCache,
) *ClientAssetsCorporation {
	return &ClientAssetsCorporation{
		naiveclient.NewNaivePagesClient[
			EntryAssetsCorporation,
			UrlParamsAssetsCorporation,
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

func NewFetchParamsAssetsCorporation(
	corporationId int32,
	refreshToken string,
) *naiveclient.NaiveClientFetchParams[UrlParamsAssetsCorporation] {
	return naiveclient.NewNaiveClientFetchParams[UrlParamsAssetsCorporation](
		UrlParamsAssetsCorporation{corporationId},
		&refreshToken,
		nil,
	)
}

// type ModelAssetsCorporation = []EntryAssetsCorporation

type EntryAssetsCorporation struct {
	// IsBlueprintCopy *bool  `json:"is_blueprint_copy,omitempty"`
	// IsSingleton     bool   `json:"is_singleton"`
	ItemId       int64  `json:"item_id"`
	LocationFlag string `json:"location_flag"`
	LocationId   int64  `json:"location_id"`
	// LocationType    string `json:"location_type"`
	Quantity int32 `json:"quantity"`
	TypeId   int32 `json:"type_id"`
}

type UrlParamsAssetsCorporation struct {
	corporationId int32
}

func (p UrlParamsAssetsCorporation) PageKey(page *int32) string {
	query := fmt.Sprintf(
		"%d/assets/?datasource=%s",
		p.corporationId,
		DATASOURCE,
	)
	query = addQueryInt32(query, "page", page)
	return query
}
func (p UrlParamsAssetsCorporation) PageUrl(page *int32) string {
	return fmt.Sprintf("%s/corporations/%s", BASE_URL, p.PageKey(page))
}
func (p UrlParamsAssetsCorporation) Key() string {
	return p.PageKey(nil)
}
func (p UrlParamsAssetsCorporation) Url() string {
	return p.PageUrl(nil)
}
func (UrlParamsAssetsCorporation) Method() string {
	return http.MethodGet
}
