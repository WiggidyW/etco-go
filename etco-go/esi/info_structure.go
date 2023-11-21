package esi

import (
	"fmt"
	"net/http"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	built "github.com/WiggidyW/etco-go/builtinconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
)

const (
	STRUCTURE_INFO_BUF_CAP        int           = 0
	STRUCTURE_INFO_MIN_EXPIRES_IN time.Duration = 24 * time.Hour
	STRUCTURE_INFO_METHOD         string        = http.MethodGet
)

type StructureInfo struct {
	Forbidden     bool   `json:"-"`
	Name          string `json:"name"`
	SolarSystemId int32  `json:"solar_system_id"`
	// OwnerId       int32             `json:"owner_id"`
	// Position      StructurePosition `json:"position"`
	// TypeId        *int32            `json:"type_id,omitempty"`
}

// type StructurePosition struct {
// 	X float64 `json:"x"`
// 	Y float64 `json:"y"`
// 	Z float64 `json:"z"`
// }

func structureInfoUrl(structureId int64) string {
	return fmt.Sprintf(
		"%s/universe/structures/%d/?datasource=%s",
		BASE_URL,
		structureId,
		DATASOURCE,
	)
}

func GetStructureInfo(x cache.Context, structureId int64) (
	rep *StructureInfo,
	expires time.Time,
	err error,
) {
	if build.STRUCTURE_INFO_WEB_REFRESH_TOKEN == built.BOOTSTRAP_STR {
		return nil, time.Now(), nil
	}
	return infoGet(
		x,
		structureInfoUrl(structureId),
		STRUCTURE_INFO_METHOD,
		keys.CacheKeyStructureInfo(structureId),
		keys.TypeStrStructureInfo,
		STRUCTURE_INFO_MIN_EXPIRES_IN,
		EsiAuthStructureInfo,
		structureInfoHandleErr,
	)
}
