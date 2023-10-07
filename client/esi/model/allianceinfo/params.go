package allianceinfo

import (
	"fmt"
	"net/http"

	"github.com/WiggidyW/etco-go/client/cachekeys"
	"github.com/WiggidyW/etco-go/client/esi/model"
)

type AllianceInfoParams struct {
	AllianceId int32
}

func (p AllianceInfoParams) CacheKey() string {
	return cachekeys.AllianceInfoCacheKey(p.AllianceId)
}

type AllianceInfoUrlParams struct {
	AllianceId int32
}

func (p AllianceInfoUrlParams) Url() string {
	return fmt.Sprintf(
		"%s/alliances/%d/?datasource=%s",
		model.BASE_URL,
		p.AllianceId,
		model.DATASOURCE,
	)
}

func (AllianceInfoUrlParams) Method() string {
	return http.MethodGet
}
