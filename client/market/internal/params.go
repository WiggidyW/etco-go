package internal

import (
	"github.com/WiggidyW/eve-trading-co-go/staticdb"
)

type MarketPriceParams struct {
	PricingInfo staticdb.PricingInfo
	TypeId      int32
}
