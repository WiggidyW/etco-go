package structureinfo

import (
	"fmt"
	"net/http"

	"github.com/WiggidyW/eve-trading-co-go/client/cachekeys"
	"github.com/WiggidyW/eve-trading-co-go/client/esi/model"
)

type StructureInfoParams struct {
	WebRefreshToken string
	StructureId     int64
}

func (p StructureInfoParams) CacheKey() string {
	return cachekeys.StructureInfoCacheKey(p.StructureId)
}

type StructureInfoUrlParams struct {
	StructureId int64
}

func (p StructureInfoUrlParams) Url() string {
	return fmt.Sprintf(
		"%s/universe/structures/%d/?datasource=%s",
		model.BASE_URL,
		p.StructureId,
		model.DATASOURCE,
	)
}

func (StructureInfoUrlParams) Method() string {
	return http.MethodGet
}
