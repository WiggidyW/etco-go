package market

import "github.com/WiggidyW/etco-go/staticdb"

type BuybackPriceParams struct {
	BuybackSystemInfo staticdb.BuybackSystemInfo
	TypeId            int32
	Quantity          int64
}
