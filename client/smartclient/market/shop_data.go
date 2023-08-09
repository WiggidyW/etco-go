package market

type ShopPrice struct {
	price float64 // price per 1 item
	desc  string
}

func newShopPrice(
	price float64,
	desc string,
) ShopPrice {
	return ShopPrice{price, desc}
}

func (s ShopPrice) Price() float64 {
	return s.price
}

func (s ShopPrice) Desc() string {
	return s.desc
}
