package market

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/staticdb"
)

func GetShopPrice(
	x cache.Context,
	typeId int32,
	quantity int64,
	locationInfo staticdb.ShopLocationInfo,
) (
	price ShopPrice,
	expires time.Time,
	err error,
) {
	return shopPriceGet(x, typeId, quantity, locationInfo)
}
