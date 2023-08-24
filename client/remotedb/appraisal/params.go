package appraisal

import (
	"github.com/WiggidyW/eve-trading-co-go/client/cachekeys"
)

type ReadAppraisalParams struct {
	AppraisalCode string
}

func (p ReadAppraisalParams) CacheKey() string {
	return cachekeys.ReadAppraisalCacheKey(p.AppraisalCode)
}
