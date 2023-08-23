package shop

func NewRejectedItem(typeId int32, quantity int64) *ShopPrice {
	return newRejected(typeId, quantity)
}
