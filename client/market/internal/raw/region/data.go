package region

import (
	"github.com/WiggidyW/weve-esi/client/market/internal/raw"
)

// int64 - locationID
type RegionMarket = map[int64][]raw.MarketOrder // not sorted or deduplicated

// faster for initialization
type initRegionMarket = map[int64]*[]raw.MarketOrder

func finishRegionMarket(init initRegionMarket) RegionMarket {
	regionMarket := make(RegionMarket, len(init))
	for locationID, orders := range init {
		regionMarket[locationID] = *orders
	}
	return regionMarket
}
