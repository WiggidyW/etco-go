package anonymous

import (
	"github.com/WiggidyW/weve-esi/client/cachekeys"
	a "github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
)

type WriteBuybackCharacterAppraisalParams[
	B a.IBuybackAppraisal[I],
	I a.IBuybackParentItem[CI],
	CI a.IBuybackChildItem,
] struct {
	CharacterId   int32
	AppraisalCode string
	IAppraisal    B
}

func (p WriteBuybackCharacterAppraisalParams[B, I, CI]) AntiCacheKey() string {
	return cachekeys.ReadCharacterAppraisalCodesCacheKey(p.CharacterId)
}
