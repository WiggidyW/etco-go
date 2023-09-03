package market

import (
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

func NewRejectedBuybackItem(
	typeId int32,
	quantity int64,
) *rdb.BuybackParentItem {
	return newRejectedParent(typeId, quantity)
}

func NewRejectedShopItem(
	typeId int32,
	quantity int64,
) *rdb.ShopItem {
	return newRejected(typeId, quantity)
}

// func minPriced(price float64, min float64) float64 {
// 	if price < min {
// 		return min
// 	}
// 	return price
// }

// func multedByModifier(price, mod float64) float64 {
// 	return mod * price / 100
// }

// func roundedToCents(price float64) float64 {
// 	return math.Round(price*100) / 100
// }
