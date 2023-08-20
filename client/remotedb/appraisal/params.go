package appraisal

import (
	"github.com/WiggidyW/weve-esi/client/cachekeys"
)

type ReadAppraisalParams struct {
	AppraisalCode string
}

func (p ReadAppraisalParams) CacheKey() string {
	return cachekeys.ReadAppraisalCacheKey(p.AppraisalCode)
}
