package writebuyback

import (
	a "github.com/WiggidyW/weve-esi/client/appraisal"
	"github.com/WiggidyW/weve-esi/client/cachekeys"
)

type WriteBuybackAppraisalParams struct {
	Appraisal a.BuybackAppraisal
}

func (p WriteBuybackAppraisalParams) AntiCacheKey() string {
	if p.Appraisal.CharacterId == nil {
		return "WBA_NULL"
	} else {
		return cachekeys.ReadUserDataCacheKey(
			*p.Appraisal.CharacterId,
		)
	}
}
