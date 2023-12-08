package implrdb

import (
	"time"
)

type UserData struct {
	BuybackAppraisals []string   `firestore:"buyback_appraisals"` // newest first, oldest last
	ShopAppraisals    []string   `firestore:"shop_appraisals"`    // newest first, oldest last
	HaulAppraisals    []string   `firestore:"haul_appraisals"`    // newest first, oldest last
	CancelledPurchase *time.Time `firestore:"cancelled_purchase"`
	MadePurchase      *time.Time `firestore:"made_purchase"`
}

// helpful when codes are stored in oldest-first order, to invert them to newest-first
func (ud *UserData) InvertCodes() {
	ud.BuybackAppraisals = inverted(ud.BuybackAppraisals)
	ud.ShopAppraisals = inverted(ud.ShopAppraisals)
	ud.HaulAppraisals = inverted(ud.HaulAppraisals)
}

func inverted[T any](slice []T) []T {
	inverted := make([]T, len(slice))
	for i, t := range slice {
		inverted[len(slice)-i-1] = t
	}
	return inverted
}
