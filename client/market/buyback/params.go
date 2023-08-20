package buyback

import "github.com/WiggidyW/weve-esi/staticdb"

type BuybackPriceParams struct {
	BuybackSystemInfo staticdb.BuybackSystemInfo
	TypeId            int32
}
