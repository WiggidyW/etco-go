package remotedb

import (
	"github.com/WiggidyW/etco-go/client/cachekeys"
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

type WriteBuybackAppraisalParams struct {
	Appraisal rdb.BuybackAppraisal
}

func (p WriteBuybackAppraisalParams) AntiCacheKey() string {
	return cachekeys.WriteBuybackAppraisalAntiCacheKey(
		p.Appraisal.CharacterId,
	)
}
