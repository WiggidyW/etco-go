package raworders_

// int64 - locationID
type RegionMarket = map[int64][]MarketOrder // not sorted or deduplicated

// faster for initialization
type initRegionMarket = map[int64]*[]MarketOrder

func finishRegionMarket(init initRegionMarket) RegionMarket {
	regionMarket := make(RegionMarket, len(init))
	for locationID, orders := range init {
		regionMarket[locationID] = *orders
	}
	return regionMarket
}
