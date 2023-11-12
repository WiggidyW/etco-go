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
	ASSETS_ENTRIES_METHOD   string = http.MethodGet
	ASSETS_ENTRIES_PER_PAGE int    = 1000
)

type AssetsEntry struct {
	ItemId       int64  `json:"item_id"`
	LocationFlag string `json:"location_flag"`
	LocationId   int64  `json:"location_id"`
	Quantity     int32  `json:"quantity"`
	TypeId       int32  `json:"type_id"`
	// IsBlueprintCopy *bool  `json:"is_blueprint_copy"`
	// IsSingleton     bool   `json:"is_singleton"`
	// LocationType    string `json:"location_type"`
}

var assetsEntriesUrl string = fmt.Sprintf(
	"%s/corporations/%d/assets/?datasource=%s",
	BASE_URL,
	build.CORPORATION_ID,
	DATASOURCE,
)

func GetAssetsEntries(x cache.Context) (
	repOrStream RepOrStream[AssetsEntry],
	expires time.Time,
	pages int,
	err error,
) {
	if build.CORPORATION_WEB_REFRESH_TOKEN == built.BOOTSTRAP_STR {
		return newBootstrapRepOrStream[AssetsEntry](), time.Now(), 0, nil
	}
	return streamGet[AssetsEntry](
		x,
		assetsEntriesUrl,
		ASSETS_ENTRIES_METHOD,
		ASSETS_ENTRIES_PER_PAGE,
		EsiAuthCorp,
	)
}
