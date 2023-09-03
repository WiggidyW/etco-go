package marketprice

import (
	"github.com/WiggidyW/etco-go/staticdb"
)

type MarketPriceParams struct {
	PricingInfo staticdb.PricingInfo
	TypeId      int32
}
