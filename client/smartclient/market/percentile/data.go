package percentile

type MrktPrctile struct {
	Price    float64 // price per 1 item
	Rejected string
}

func newMrktPrctile(price float64, rejected string) MrktPrctile {
	return MrktPrctile{price, rejected}
}
