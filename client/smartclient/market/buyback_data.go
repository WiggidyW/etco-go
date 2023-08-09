package market

type BuybackPrice struct {
	price    float64 // price per 1 item
	desc     string
	Children []BuybackPriceChild
}

func newBuybackPrice(
	price float64,
	desc string,
	Children []BuybackPriceChild,
) BuybackPrice {
	return BuybackPrice{price, desc, Children}
}

func (b BuybackPrice) Price() float64 {
	return b.price
}

func (b BuybackPrice) Desc() string {
	return b.desc
}

type BuybackPriceChild struct {
	price    float64 // price per 1 item
	desc     string
	Quantity float64 // number per 1 parent item
}

func newBuybackPriceChild(
	price float64,
	desc string,
	quantity float64,
) BuybackPriceChild {
	return BuybackPriceChild{price, desc, quantity}
}

func (bc BuybackPriceChild) Price() float64 {
	return bc.price
}

func (bc BuybackPriceChild) Desc() string {
	return bc.desc
}
