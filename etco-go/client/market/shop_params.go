package market

import "github.com/WiggidyW/etco-go/staticdb"

type ShopPriceParams struct {
	ShopLocationInfo staticdb.ShopLocationInfo
	TypeId           int32
	Quantity         int64
}
