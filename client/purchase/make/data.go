package make

import "github.com/WiggidyW/weve-esi/client/appraisal"

type MakePurchaseStatus uint8

const (
	Success        MakePurchaseStatus = iota
	CooldownLimit                     // appraisal is nil
	MaxActiveLimit                    // appraisal is nil
	ItemsRejectedAndUnavailable
	ItemsRejected
	ItemsUnavailable
)

type MakePurchaseResponse struct {
	Status    MakePurchaseStatus
	Appraisal *appraisal.ShopAppraisal
}
