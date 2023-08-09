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
	CONTRACT_ITEMS_CORPORATION_USE_AUTH         bool  = true
	CONTRACT_ITEMS_CORPORATION_ENTRIES_PER_PAGE int32 = 5000
)

type ClientContractItemsCorporation struct {
	naiveclient.CachingNaivePageEntriesClient[
		EntryContractItemsCorporation,
		UrlParamsContractItemsCorporation,
	]
}

func NewClientContractItemsCorporation(
	rawClient *rawclient.RawClient,
	minExpires time.Duration,
	modelServerLockTTL time.Duration,
	modelServerLockMaxWait time.Duration,
	clientCache cache.SharedClientCache,
	serverCache cache.SharedServerCache,
) *ClientContractItemsCorporation {
	cachingClient := naiveclient.NewCachingNaivePageEntriesClient[
		EntryContractItemsCorporation,
		UrlParamsContractItemsCorporation,
	](
		rawClient,
		CONTRACT_ITEMS_CORPORATION_USE_AUTH,
		CONTRACT_ITEMS_CORPORATION_ENTRIES_PER_PAGE,
		minExpires,
		cache.NewBufferPool(0),
		clientCache,
		serverCache,
		modelServerLockTTL,
		modelServerLockMaxWait,
	)
	return &ClientContractItemsCorporation{cachingClient}
}

func NewFetchParamsContractItemsCorporation(
	corporationId int32,
	contractId int32,
	refreshToken string,
) *naiveclient.NaiveClientFetchParams[UrlParamsContractItemsCorporation] {
	return naiveclient.NewNaiveClientFetchParams[UrlParamsContractItemsCorporation](
		UrlParamsContractItemsCorporation{corporationId, contractId},
		&refreshToken,
		nil,
	)
}

// type ModelContractItemsCorporation = []EntryContractItemsCorporation

type EntryContractItemsCorporation struct {
	IsIncluded  bool   `json:"is_included"`
	IsSingleton bool   `json:"is_singleton"`
	Quantity    int32  `json:"quantity"`
	RawQuantity *int32 `json:"raw_quantity,omitempty"`
	RecordId    int64  `json:"record_id"`
	TypeId      int32  `json:"type_id"`
}

type UrlParamsContractItemsCorporation struct {
	corporationId int32
	contractId    int32
}

func (p UrlParamsContractItemsCorporation) CacheKey() string {
	return fmt.Sprintf(
		"%d/contracts/%d/items/?datasource=%s",
		p.corporationId,
		p.contractId,
		DATASOURCE,
	)
}
func (p UrlParamsContractItemsCorporation) Url() string {
	return fmt.Sprintf("%s/corporations/%s", BASE_URL, p.CacheKey())
}
func (UrlParamsContractItemsCorporation) Method() string {
	return http.MethodGet
}
