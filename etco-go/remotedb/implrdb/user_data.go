package implrdb

import (
	"time"
)

type UserData struct {
	BuybackAppraisals []string   `firestore:"buyback_appraisals"`
	ShopAppraisals    []string   `firestore:"shop_appraisals"`
	CancelledPurchase *time.Time `firestore:"cancelled_purchase"`
	MadePurchase      *time.Time `firestore:"made_purchase"`
}
