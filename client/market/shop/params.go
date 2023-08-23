package shop

import "github.com/WiggidyW/weve-esi/staticdb"

type ShopPriceParams struct {
	ShopLocationInfo staticdb.ShopLocationInfo
	TypeId           int32
	Quantity         int64
}
