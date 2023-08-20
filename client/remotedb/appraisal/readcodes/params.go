package readcodes

import (
	"github.com/WiggidyW/weve-esi/client/cachekeys"
)

type ReadCharacterAppraisalCodesParams struct {
	CharacterId int32
}

func (p ReadCharacterAppraisalCodesParams) CacheKey() string {
	return cachekeys.ReadCharacterAppraisalCodesCacheKey(p.CharacterId)
}
