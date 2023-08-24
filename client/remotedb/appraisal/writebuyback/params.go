package writebuyback

import (
	a "github.com/WiggidyW/eve-trading-co-go/client/appraisal"
	"github.com/WiggidyW/eve-trading-co-go/client/cachekeys"
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
