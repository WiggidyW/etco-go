package userdata

import (
	"time"

	"github.com/WiggidyW/etco-go/client/contracts"
)

type BuybackAppraisalStatus struct {
	Code     string
	Contract *contracts.Contract
}

type ShopAppraisalStatus struct {
	Code            string
	Contract        *contracts.Contract
	InPurchaseQueue bool
}

type UserData struct {
	BuybackAppraisals []BuybackAppraisalStatus // oldest first
	ShopAppraisals    []ShopAppraisalStatus    // oldest first
	CancelledPurchase time.Time
	MadePurchase      time.Time
}
