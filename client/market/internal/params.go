package internal

import (
	"github.com/WiggidyW/weve-esi/staticdb"
)

type MarketPriceParams struct {
	PricingInfo staticdb.PricingInfo
	TypeId      int32
}
