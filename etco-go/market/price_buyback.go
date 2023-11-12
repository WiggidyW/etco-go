package market

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/staticdb"
)

func GetBuybackPrice(
	x cache.Context,
	typeId int32,
	quantity int64,
	systemInfo staticdb.BuybackSystemInfo,
) (
	price BuybackPriceParent,
	expires time.Time,
	err error,
) {
	return buybackPriceGet(x, typeId, quantity, systemInfo)
}
