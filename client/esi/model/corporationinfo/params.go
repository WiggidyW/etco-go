package corporationinfo

import (
	"fmt"
	"net/http"

	"github.com/WiggidyW/etco-go/client/cachekeys"
	"github.com/WiggidyW/etco-go/client/esi/model"
)

type CorporationInfoParams struct {
	CorporationId int32
}

func (p CorporationInfoParams) CacheKey() string {
	return cachekeys.CorporationInfoCacheKey(p.CorporationId)
}

type CorporationInfoUrlParams struct {
	CorporationId int32
}

func (p CorporationInfoUrlParams) Url() string {
	return fmt.Sprintf(
		"%s/corporations/%d/?datasource=%s",
		model.BASE_URL,
		p.CorporationId,
		model.DATASOURCE,
	)
}

func (CorporationInfoUrlParams) Method() string {
	return http.MethodGet
}
