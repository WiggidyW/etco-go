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
	CONTRACT_ITEMS_ENTRIES_MAKE_CAP       int           = 100
	CONTRACT_ITEMS_ENTRIES_MIN_EXPIRES_IN time.Duration = 24 * time.Hour
	CONTRACT_ITEMS_ENTRIES_METHOD         string        = http.MethodGet
)

type ContractItemsEntry struct {
	Quantity int32 `json:"quantity"`
	TypeId   int32 `json:"type_id"`
	// IsIncluded  bool   `json:"is_included"`
	// IsSingleton bool   `json:"is_singleton"`
	// RawQuantity *int32 `json:"raw_quantity"`
	// RecordId    int64  `json:"record_id"`
}

func contractItemsEntriesUrl(contractId int32) string {
	return fmt.Sprintf(
		"%s/corporations/%d/contracts/%d/items/?datasource=%s",
		BASE_URL,
		build.CORPORATION_ID,
		contractId,
		DATASOURCE,
	)
}

func GetContractItemsEntries(x cache.Context, contractId int32) (
	rep []ContractItemsEntry,
	expires time.Time,
	err error,
) {
	if build.CORPORATION_WEB_REFRESH_TOKEN == built.BOOTSTRAP_STR {
		return nil, expires, nil
	}
	return contractItemsEntriesGet(x, contractId)
}
