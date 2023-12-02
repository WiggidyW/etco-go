package esi

import (
	"fmt"
	"net/http"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	built "github.com/WiggidyW/etco-go/builtinconstants"
	"github.com/WiggidyW/etco-go/cache"
)

const (
	CONTRACTS_ENTRIES_METHOD   string = http.MethodGet
	CONTRACTS_ENTRIES_PER_PAGE int    = 1000
)

type ContractsEntry struct {
	AssigneeId          int32     `json:"assignee_id"`
	Availability        string    `json:"availability"`
	ContractId          int32     `json:"contract_id"`
	DateExpired         time.Time `json:"date_expired"`
	DateIssued          time.Time `json:"date_issued"`
	StartLocationId     *int64    `json:"start_location_id"`
	EndLocationId       *int64    `json:"end_location_id"`
	IssuerCorporationId int32     `json:"issuer_corporation_id"`
	IssuerId            int32     `json:"issuer_id"`
	Price               *float64  `json:"price"`
	Reward              *float64  `json:"reward"`
	Status              string    `json:"status"`
	Title               *string   `json:"title,omitempty"`
	Type                string    `json:"type"`
	Volume              *float64  `json:"volume,omitempty"`
	Collateral          *float64  `json:"collateral"`
	// AcceptorId   int32  `json:"acceptor_id"`
	// Buyout              *float64   `json:"buyout"`
	// DateAccepted        *time.Time `json:"date_accepted"`
	// DateCompleted       *time.Time `json:"date_completed"`
	// DaysToComplete      *int32     `json:"days_to_complete"`
	// ForCorporation      bool       `json:"for_corporation"`
}

var contractsEntriesUrl string = fmt.Sprintf(
	"%s/corporations/%d/contracts/?datasource=%s",
	BASE_URL,
	build.CORPORATION_ID,
	DATASOURCE,
)

func GetContractsEntries(x cache.Context) (
	repOrStream RepOrStream[ContractsEntry],
	expires time.Time,
	pages int,
	err error,
) {
	if build.CORPORATION_WEB_REFRESH_TOKEN == built.BOOTSTRAP_STR {
		return newBootstrapRepOrStream[ContractsEntry](), time.Now(), 0, nil
	}
	return streamGet[ContractsEntry](
		x,
		contractsEntriesUrl,
		CONTRACTS_ENTRIES_METHOD,
		CONTRACTS_ENTRIES_PER_PAGE,
		EsiAuthCorp,
	)
}
