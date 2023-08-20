package anonymous

import (
	a "github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
)

type WriteBuybackAnonAppraisalParams[
	B a.IBuybackAppraisal[I],
	I a.IBuybackParentItem[CI],
	CI a.IBuybackChildItem,
] struct {
	AppraisalCode string
	IAppraisal    B
}
