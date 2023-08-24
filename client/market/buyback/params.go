package buyback

import "github.com/WiggidyW/eve-trading-co-go/staticdb"

type BuybackPriceParams struct {
	BuybackSystemInfo staticdb.BuybackSystemInfo
	TypeId            int32
	Quantity          int64
}
