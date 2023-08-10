package percentile

import (
	"context"
	"fmt"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/smartclient/market/desc"
	"github.com/WiggidyW/weve-esi/client/smartclient/market/percentile/orders"
	"github.com/WiggidyW/weve-esi/logger"
)

type MrktPrctileParams struct {
	orders.MrktOrdersParams
}

func (f MrktPrctileParams) CacheKey() string {
	return fmt.Sprintf(
		"mrktprice-%s-%d",
		f.CacheKeyInner(),
		f.PricingInfo.Prctile,
	)
}

type MrktPrctileClient struct {
	client *client.CachingClient[
		orders.MrktOrdersParams,
		orders.MrktOrders,
		cache.ExpirableData[orders.MrktOrders],
		*orders.MrktOrdersClient,
	]
}

func (mpc *MrktPrctileClient) Fetch(
	ctx context.Context,
	params MrktPrctileParams,
) (*cache.ExpirableData[MrktPrctile], error) {
	// fetch the mrkt orders
	mrktOrders, err := mpc.client.Fetch(
		ctx,
		params.MrktOrdersParams,
	)
	if err != nil {
		return nil, err
	}
	expires := mrktOrders.Expires()

	// if no orders, return rejected
	if !mrktOrders.Data().HasOrders() {
		data := cache.NewExpirableData[MrktPrctile](
			newMrktPrctile(0, desc.RejectedNoOrders(
				params.PricingInfo.MrktName,
			)),
			expires,
		)
		return &data, nil
	}

	// get the prctile price
	price, ok := mrktOrders.Data().Prctile(params.PricingInfo.Prctile)
	if !ok {
		// warn if it isn't pre-calculated
		logger.Logger.Warn(fmt.Sprintf(
			"prctile %d not pre-calculated for %d at %s",
			params.PricingInfo.Prctile,
			params.TypeId,
			params.PricingInfo.MrktName,
		))
	}

	// validate the prctile price
	if price <= 0 {
		// if it's <= 0, warn and return rejected
		logger.Logger.Warn(fmt.Sprintf(
			"prctile %d had price %f for %d at %s, but orders exist",
			params.PricingInfo.Prctile,
			price,
			params.TypeId,
			params.PricingInfo.MrktName,
		))
		data := cache.NewExpirableData[MrktPrctile](
			newMrktPrctile(0, desc.RejectedServerError()),
			expires,
		)
		return &data, nil
	}

	// return the prctile price
	data := cache.NewExpirableData[MrktPrctile](
		newMrktPrctile(price, ""),
		expires,
	)
	return &data, nil
}
