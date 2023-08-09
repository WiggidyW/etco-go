package percentile

import (
	"context"
	"fmt"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/smartclient/market/desc"
	"github.com/WiggidyW/weve-esi/client/smartclient/market/percentile/orders"
	"github.com/WiggidyW/weve-esi/logger"
	"github.com/WiggidyW/weve-esi/staticdb/tc"
)

type MarketPercentileClientFetchParams struct {
	orders.MarketOrdersClientFetchParams
}

func NewFetchParams(
	pricingInfo *tc.PricingInfo,
	typeId int32,
) MarketPercentileClientFetchParams {
	return MarketPercentileClientFetchParams{
		orders.NewFetchParams(pricingInfo, typeId),
	}
}

func (f MarketPercentileClientFetchParams) CacheKey() string {
	return fmt.Sprintf(
		"marketprice-%s-%d",
		f.CacheKeyInner(),
		f.PricingInfo.Percentile(),
	)
}

type MarketPercentileClient struct {
	client *client.CachingClient[
		orders.MarketOrdersClientFetchParams,
		orders.MarketOrders,
		cache.ExpirableData[orders.MarketOrders],
		*orders.MarketOrdersClient,
	]
}

func (mpc *MarketPercentileClient) Fetch(
	ctx context.Context,
	params MarketPercentileClientFetchParams,
) (*cache.ExpirableData[MarketPercentile], error) {
	// fetch the market orders
	marketOrders, err := mpc.client.Fetch(
		ctx,
		params.MarketOrdersClientFetchParams,
	)
	if err != nil {
		return nil, err
	}
	expires := marketOrders.Expires()

	// if no orders, return rejected
	if !marketOrders.Data().HasOrders() {
		data := cache.NewExpirableData[MarketPercentile](
			newMarketPercentile(0, desc.RejectedNoOrders(
				params.PricingInfo.MarketName(),
			)),
			expires,
		)
		return &data, nil
	}

	// validate the given percentile
	percentile := params.PricingInfo.Percentile()
	if percentile < 0 || percentile > 100 {
		// if it's invalid, warn and return rejected
		logger.Logger.Warn(fmt.Sprintf(
			"percentile %d not between 0 and 100 for %d at %s",
			percentile,
			params.TypeId,
			params.PricingInfo.MarketName(),
		))
		data := cache.NewExpirableData[MarketPercentile](
			newMarketPercentile(0, desc.RejectedServerError()),
			expires,
		)
		return &data, nil
	}

	// get the percentile price
	price, ok := marketOrders.Data().Percentile(percentile)
	if !ok {
		// warn if it isn't pre-calculated
		logger.Logger.Warn(fmt.Sprintf(
			"percentile %d not pre-calculated for %d at %s",
			percentile,
			params.TypeId,
			params.PricingInfo.MarketName(),
		))
	}

	// validate the percentile price
	if price <= 0 {
		// if it's <= 0, warn and return rejected
		logger.Logger.Warn(fmt.Sprintf(
			"percentile %d had price %f for %d at %s, but orders exist",
			percentile,
			price,
			params.TypeId,
			params.PricingInfo.MarketName(),
		))
		data := cache.NewExpirableData[MarketPercentile](
			newMarketPercentile(0, desc.RejectedServerError()),
			expires,
		)
		return &data, nil
	}

	// return the percentile price
	data := cache.NewExpirableData[MarketPercentile](
		newMarketPercentile(price, ""),
		expires,
	)
	return &data, nil
}
