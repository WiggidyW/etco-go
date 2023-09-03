package remotedb

import (
	"github.com/WiggidyW/etco-go/client/cachekeys"
)

type ReadAppraisalParams struct {
	AppraisalCode string
}

func (p ReadAppraisalParams) CacheKey() string {
	return cachekeys.ReadAppraisalCacheKey(p.AppraisalCode)
}
