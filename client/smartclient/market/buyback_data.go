package market

type BuybackPrice struct {
	price    float64 // price per 1 item
	desc     string
	Children []BuybackPriceChild
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

func (bc BuybackPriceChild) Price() float64 {
	return bc.price
}

func (bc BuybackPriceChild) Desc() string {
	return bc.desc
}
