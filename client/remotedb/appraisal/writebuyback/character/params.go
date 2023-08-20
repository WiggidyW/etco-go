package anonymous

import (
	a "github.com/WiggidyW/weve-esi/client/appraisal"
	"github.com/WiggidyW/weve-esi/client/cachekeys"
)

type WriteBuybackCharacterAppraisalParams struct {
	CharacterId   int32
	AppraisalCode string
	Appraisal     a.BuybackAppraisal
}

func (p WriteBuybackCharacterAppraisalParams) AntiCacheKey() string {
	return cachekeys.ReadCharacterAppraisalCodesCacheKey(p.CharacterId)
}
