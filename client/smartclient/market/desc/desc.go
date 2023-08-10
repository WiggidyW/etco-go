package desc

import "fmt"

const (
	_REJECTED              string = "Rejected"
	_REJECTED_SERVER_ERROR string = "Rejected - server error"
)

func Rejected() string {
	return _REJECTED
}

func RejectedServerError() string {
	return _REJECTED_SERVER_ERROR
}

func RejectedNoOrders(market string) string {
	return "Rejected - no orders found at " + market
}

func Accepted(
	market string,
	percentile int,
	modifier float64,
	isBuy bool,
) string {
	var percentileStr string
	if percentile == 0 {
		if isBuy {
			percentileStr = "MinBuy"
		} else {
			percentileStr = "MinSell"
		}
	} else if percentile == 100 {
		if isBuy {
			percentileStr = "MaxBuy"
		} else {
			percentileStr = "MaxSell"
		}
	} else {
		if isBuy {
			percentileStr = fmt.Sprintf(
				"%dth Percentile Buy",
				percentile,
			)
		} else {
			percentileStr = fmt.Sprintf(
				"%dth Percentile Sell",
				percentile,
			)
		}
	}
	return fmt.Sprintf(
		"%s %d%% of %s",
		market,
		uint8(modifier*100),
		percentileStr,
	)
}

func AcceptedWithFee(
	market string,
	percentile int,
	modifier float64,
	isBuy bool,
	fee float64,
) string {
	return fmt.Sprintf(
		"%s (fee: %.2f)",
		Accepted(market, percentile, modifier, isBuy),
		fee,
	)
}

func AcceptedReprocessed(repEff float64) string {
	return fmt.Sprintf(
		"%.0f%% Reprocessed",
		repEff,
	)
}

func RejectedReprocessed(repEff float64) string {
	return "Rejected - " +
		AcceptedReprocessed(repEff) +
		" - All reprocessed items rejected"
}
