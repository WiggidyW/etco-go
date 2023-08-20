package market

import "fmt"

const (
	_REJECTED              string = "Rejected"
	_REJECTED_SERVER_ERROR string = "Rejected - server error"
	// we might have other fees than m3 at some point, but for now
	_REJECTED_FEE string = "Rejected - m3 fee makes price negative"
)

func Rejected() string {
	return _REJECTED
}

func RejectedServerError() string {
	return _REJECTED_SERVER_ERROR
}

func RejectedFee() string {
	return _REJECTED_FEE
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
		var percentileValue string
		if percentile >= 4 {
			percentileValue = fmt.Sprintf("%dth", percentile)
		} else if percentile == 3 {
			percentileValue = "3rd"
		} else if percentile == 2 {
			percentileValue = "2nd"
		} else if percentile == 1 {
			percentileValue = "1st"
		}
		if isBuy {
			percentileStr = fmt.Sprintf(
				"%s Percentile Buy",
				percentileValue,
			)
		} else {
			percentileStr = fmt.Sprintf(
				"%s Percentile Sell",
				percentileValue,
			)
		}
	}
	return fmt.Sprintf(
		"%s %.0f%% of %s",
		market,
		modifier,
		percentileStr,
	)
}

func AcceptedReprocessed(repEff float64) string {
	return fmt.Sprintf(
		"%.0f%% Reprocessed",
		repEff,
	)
}

func RejectedReprocessed(repEff float64) string {
	return fmt.Sprintf(
		"Rejected - %s - all reprocessed items rejected",
		AcceptedReprocessed(repEff),
	)
}

func RejectedReprocessedFee(repEff float64) string {
	return fmt.Sprintf(
		"Rejected - %s - m3 fee makes price negative",
		AcceptedReprocessed(repEff),
	)
}
