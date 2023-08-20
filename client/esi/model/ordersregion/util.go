package ordersregion

var (
	BUY  string = "buy"
	SELL string = "sell"
	// ALL string = "all"
)

func boolToOrderType(b bool) *string {
	if b {
		return &BUY
	} else {
		return &SELL
	}
}
