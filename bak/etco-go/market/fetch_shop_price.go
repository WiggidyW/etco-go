package market

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/staticdb"
)

func shopPriceGet(
	x cache.Context,
	typeId int32,
	quantity int64,
	shopLocationInfo staticdb.ShopLocationInfo,
) (
	price ShopPrice,
	expires time.Time,
	err error,
) {
	return fetch.HandleFetch(
		x,
		nil,
		shopPriceGetFetchFunc(typeId, quantity, shopLocationInfo),
		nil,
	)
}

func shopPriceGetFetchFunc(
	typeId int32,
	quantity int64,
	locationInfo staticdb.ShopLocationInfo,
) fetch.Fetch[ShopPrice] {
	return func(x cache.Context) (
		price ShopPrice,
		expires time.Time,
		_ *postfetch.Params,
		err error,
	) {
		pricingInfo := locationInfo.GetTypePricingInfo(typeId)
		if pricingInfo == nil {
			price = newRejected(typeId, quantity)
			expires = fetch.MAX_EXPIRES
		} else {
			var positivePrice float64
			positivePrice, expires, err = GetPercentilePrice(
				x,
				typeId,
				*pricingInfo,
			)
			if err != nil {
				return price, expires, nil, err
			}
			price = unpackPositivePrice(
				typeId,
				quantity,
				positivePrice,
				*pricingInfo,
			)
		}
		return price, expires, nil, nil
	}
}
