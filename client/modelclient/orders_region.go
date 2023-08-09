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
	ORDERS_REGION_USE_AUTH         bool  = false
	ORDERS_REGION_ENTRIES_PER_PAGE int32 = 1000
)

var (
	BUY  string = "buy"
	SELL string = "sell"
	// ALL string = "all"
)

type ClientOrdersRegion struct {
	naiveclient.NaivePagesClient[
		EntryOrdersRegion,
		UrlParamsOrdersRegion,
	]
}

func NewClientOrdersRegion(
	rawClient *rawclient.RawClient,
	minExpires time.Duration,
	modelServerLockTTL time.Duration,
	modelServerLockMaxWait time.Duration,
	headPool *cache.BufferPool,
	headServerLockTTL time.Duration,
	headServerLockMaxWait time.Duration,
	clientCache cache.SharedClientCache,
	serverCache cache.SharedServerCache,
) *ClientOrdersRegion {
	return &ClientOrdersRegion{
		naiveclient.NewNaivePagesClient[
			EntryOrdersRegion,
			UrlParamsOrdersRegion,
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

func NewFetchParamsOrdersRegion(
	regionId int32,
	typeId *int32,
	isBuy *bool,
) *naiveclient.NaiveClientFetchParams[UrlParamsOrdersRegion] {
	return naiveclient.NewNaiveClientFetchParams[UrlParamsOrdersRegion](
		UrlParamsOrdersRegion{regionId, typeId, boolToOrderType(isBuy)},
		nil,
		nil,
	)
}

// type ModelOrdersRegion = []EntryOrdersRegion

type EntryOrdersRegion struct {
	// Duration   int32 `json:"duration"`
	// IsBuyOrder bool  `json:"is_buy_order"`
	// Issued       time.Time `json:"issued"`
	LocationId int64 `json:"location_id"`
	// MinVolume   int32 `json:"min_volume"`
	// OrderId     int64 `json:"order_id"`
	Price float64 `json:"price"`
	// Range       string `json:"range"`
	// TypeId       int32 `json:"type_id"`
	VolumeRemain int32 `json:"volume_remain"`
	// VolumeTotal int32 `json:"volume_total"`
}

type UrlParamsOrdersRegion struct {
	regionId  int32
	typeId    *int32
	orderType *string
}

func (p UrlParamsOrdersRegion) PageKey(page *int32) string {
	query := fmt.Sprintf(
		"%d/orders/?datasource=%s",
		p.regionId,
		DATASOURCE,
	)
	query = addQueryInt32(query, "page", page)
	query = addQueryInt32(query, "type_id", p.typeId)
	query = addQueryString(query, "order_type", p.orderType)
	return query
}
func (p UrlParamsOrdersRegion) PageUrl(page *int32) string {
	return fmt.Sprintf("%s/markets/%s", BASE_URL, p.PageKey(page))
}
func (p UrlParamsOrdersRegion) Key() string {
	return p.PageKey(nil)
}
func (p UrlParamsOrdersRegion) Url() string {
	return p.PageUrl(nil)
}
func (UrlParamsOrdersRegion) Method() string {
	return http.MethodGet
}

func boolToOrderType(b *bool) *string {
	if b == nil {
		// return &ALL
		return nil
	} else if *b {
		return &BUY
	} else {
		return &SELL
	}
}
