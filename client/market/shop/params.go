package shop

import "github.com/WiggidyW/eve-trading-co-go/staticdb"

type ShopPriceParams struct {
	ShopLocationInfo staticdb.ShopLocationInfo
	TypeId           int32
	Quantity         int64
}
