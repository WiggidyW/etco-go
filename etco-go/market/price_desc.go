package market

import (
	"fmt"
	"math"
)

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
		if percentile == 3 ||
			percentile == 23 ||
			percentile == 33 ||
			percentile == 43 ||
			percentile == 53 ||
			percentile == 63 ||
			percentile == 73 ||
			percentile == 83 ||
			percentile == 93 {
			percentileValue = fmt.Sprintf("%drd", percentile)
		} else if percentile == 2 ||
			percentile == 22 ||
			percentile == 32 ||
			percentile == 42 ||
			percentile == 52 ||
			percentile == 62 ||
			percentile == 72 ||
			percentile == 82 ||
			percentile == 92 {
			percentileValue = fmt.Sprintf("%dnd", percentile)
		} else if percentile == 1 ||
			percentile == 21 ||
			percentile == 31 ||
			percentile == 41 ||
			percentile == 51 ||
			percentile == 61 ||
			percentile == 71 ||
			percentile == 81 ||
			percentile == 91 {
			percentileValue = fmt.Sprintf("%dst", percentile)
		} else {
			percentileValue = fmt.Sprintf("%dth", percentile)
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
		math.Round(modifier*100.0),
		percentileStr,
	)
}

func AcceptedReprocessed(repEff float64) string {
	return fmt.Sprintf(
		"%.0f%% Reprocessed",
		math.Round(repEff*100.0),
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
