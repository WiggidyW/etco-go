package userdata

import (
	"time"

	"github.com/WiggidyW/etco-go/contracts"
)

type AppraisalStatus struct {
	Code     string
	Contract *contracts.Contract
}

type ShopAppraisalStatus struct {
	InPurchaseQueue bool
	AppraisalStatus
}

func getAppraisalStatuses(
	codes []string,
	contracts map[string]contracts.Contract,
) (
	appraisalStatus AppraisalStatus,
	expires time.Time,
	err error,
) {

}
