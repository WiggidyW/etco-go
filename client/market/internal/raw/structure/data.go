package structure

import (
	"github.com/WiggidyW/eve-trading-co-go/client/market/internal/raw"
)

// int32 - typeID
type StructureMarket = map[int32]StructureOrders // not sorted or deduplicated

type StructureOrders struct {
	Buy  []raw.MarketOrder
	Sell []raw.MarketOrder
}

// faster for initialization
// ~0.25s vs ~0.18s including the finish function
type initStructureMarket = map[int32]*StructureOrders

func finishStructureMarket(init initStructureMarket) StructureMarket {
	structureMarket := make(StructureMarket, len(init))
	for locationID, orders := range init {
		structureMarket[locationID] = *orders
	}
	return structureMarket
}
