package buyback

func NewRejectedItem(typeId int32, quantity int64) *BuybackPriceParent {
	return newRejectedParent(typeId, quantity)
}
