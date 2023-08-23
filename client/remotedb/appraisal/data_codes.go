package appraisal

import "time"

const (
	CHARACTERS_COLLECTION_ID string = "characters"

	S_CHAR_TIME_CANCELLED_PURCHASE string = "cancelled_purchase"
	S_CHAR_TIME_MADE_PURCHASE      string = "made_purchase"
	S_CHAR_APPRAISALS              string = "shop_appraisals"
	B_CHAR_APPRAISALS              string = "buyback_appraisals"
)

type UserData struct {
	ShopAppraisals    []string  `firestore:"shop_appraisals"`
	BuybackAppraisals []string  `firestore:"buyback_appraisals"`
	CancelledPurchase time.Time `firestore:"cancelled_purchase"`
	MadePurchase      time.Time `firestore:"made_purchase"`
}
