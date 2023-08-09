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
	CONTRACTS_CORPORATION_USE_AUTH         bool  = true
	CONTRACTS_CORPORATION_ENTRIES_PER_PAGE int32 = 1000
)

type ClientContractsCorporation struct {
	naiveclient.NaivePagesClient[
		EntryContractsCorporation,
		UrlParamsContractsCorporation,
	]
}

func NewClientContractsCorporation(
	rawClient *rawclient.RawClient,
	minExpires time.Duration,
	modelServerLockTTL time.Duration,
	modelServerLockMaxWait time.Duration,
	headPool *cache.BufferPool,
	headServerLockTTL time.Duration,
	headServerLockMaxWait time.Duration,
	clientCache cache.SharedClientCache,
	serverCache cache.SharedServerCache,
) *ClientContractsCorporation {
	return &ClientContractsCorporation{
		naiveclient.NewNaivePagesClient[
			EntryContractsCorporation,
			UrlParamsContractsCorporation,
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

func NewFetchParamsContractsCorporation(
	corporationId int32,
	refreshToken string,
) *naiveclient.NaiveClientFetchParams[UrlParamsContractsCorporation] {
	return naiveclient.NewNaiveClientFetchParams[UrlParamsContractsCorporation](
		UrlParamsContractsCorporation{corporationId},
		&refreshToken,
		nil,
	)
}

// type ModelContractsCorporation = []EntryContractsCorporation

type EntryContractsCorporation struct {
	AcceptorId          int32      `json:"acceptor_id"`
	AssigneeId          int32      `json:"assignee_id"`
	Availability        string     `json:"availability"`
	Buyout              *float64   `json:"buyout,omitempty"`
	Collateral          *float64   `json:"collateral,omitempty"`
	ContractId          int32      `json:"contract_id"`
	DateAccepted        *time.Time `json:"date_accepted,omitempty"`
	DateCompleted       *time.Time `json:"date_completed,omitempty"`
	DateExpired         time.Time  `json:"date_expired"`
	DateIssued          time.Time  `json:"date_issued"`
	DaysToComplete      *int32     `json:"days_to_complete,omitempty"`
	EndLocationId       *int64     `json:"end_location_id,omitempty"`
	ForCorporation      bool       `json:"for_corporation"`
	IssuerCorporationId int32      `json:"issuer_corporation_id"`
	IssuerId            int32      `json:"issuer_id"`
	Price               *float64   `json:"price,omitempty"`
	Reward              *float64   `json:"reward,omitempty"`
	StartLocationId     *int64     `json:"start_location_id,omitempty"`
	Status              string     `json:"status"`
	Title               *string    `json:"title,omitempty"`
	Type                string     `json:"type"`
	Volume              *float64   `json:"volume,omitempty"`
}

type UrlParamsContractsCorporation struct {
	corporationId int32
}

func (p UrlParamsContractsCorporation) PageKey(page *int32) string {
	query := fmt.Sprintf(
		"%d/contracts/?datasource=%s",
		p.corporationId,
		DATASOURCE,
	)
	query = addQueryInt32(query, "page", page)
	return query
}
func (p UrlParamsContractsCorporation) PageUrl(page *int32) string {
	return fmt.Sprintf("%s/corporations/%s", BASE_URL, p.PageKey(page))
}
func (p UrlParamsContractsCorporation) Key() string {
	return p.PageKey(nil)
}
func (p UrlParamsContractsCorporation) Url() string {
	return p.PageUrl(nil)
}
func (UrlParamsContractsCorporation) Method() string {
	return http.MethodGet
}
