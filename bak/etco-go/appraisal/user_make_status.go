package appraisal

type MakePurchaseStatus uint8

const (
	MPS_Success        MakePurchaseStatus = iota
	MPS_CooldownLimit                     // appraisal is nil
	MPS_MaxActiveLimit                    // appraisal is nil
	MPS_ItemsRejectedAndUnavailable
	MPS_ItemsRejected
	MPS_ItemsUnavailable
)
