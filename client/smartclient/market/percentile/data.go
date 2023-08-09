package percentile

type MarketPercentile struct {
	Price    float64 // price per 1 item
	Rejected string
}

func newMarketPercentile(price float64, rejected string) MarketPercentile {
	return MarketPercentile{price, rejected}
}
