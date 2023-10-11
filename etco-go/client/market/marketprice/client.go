package marketprice

import (
	"context"
	"fmt"

	"github.com/WiggidyW/etco-go/cache"
	mordersregion "github.com/WiggidyW/etco-go/client/esi/model/ordersregion"
	mordersstructure "github.com/WiggidyW/etco-go/client/esi/model/ordersstructure"
	"github.com/WiggidyW/etco-go/client/market/marketprice/filterorders_"
	"github.com/WiggidyW/etco-go/logger"
)

type MarketPriceClient struct {
	structureClient filterorders_.WC_FilterStructureMarketClient
	regionClient    filterorders_.WC_FilterRegionMarketClient
}

func NewMarketPriceClient(
	modelRegionClient mordersregion.OrdersRegionClient,
	modelStructureClient mordersstructure.OrdersStructureClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) MarketPriceClient {
	return MarketPriceClient{
		filterorders_.NewWC_FilterStructureMarketClient(
			modelStructureClient,
			cCache,
			sCache,
		),
		filterorders_.NewWC_FilterRegionMarketClient(
			modelRegionClient,
			cCache,
			sCache,
		),
	}
}

// returns nil if there are no orders to be found
func (mpc MarketPriceClient) Fetch(
	ctx context.Context,
	params MarketPriceParams,
) (*PositivePrice, error) {
	var sortedOrders filterorders_.SortedMarketOrders
	var err error

	if params.PricingInfo.MarketIsStructure {
		sortedOrders, err = mpc.fetchStructure(ctx, params)
	} else {
		sortedOrders, err = mpc.fetchStation(ctx, params)
	}
	if err != nil {
		return nil, err
	} else if sortedOrders.Quantity == 0 {
		return nil, nil
	}

	price := sortedOrders.PercentilePrice(params.PricingInfo.Percentile)
	price = price * params.PricingInfo.Modifier

	if price <= 0 {
		logger.Warn(fmt.Errorf(
			"price for type '%d' with pricing info '%+v' was '%f'",
			params.TypeId,
			params.PricingInfo,
			price,
		))
		return nil, nil
	}

	positivePrice := PositivePrice(price)
	return &positivePrice, nil
}

func (mpc MarketPriceClient) fetchStructure(
	ctx context.Context,
	params MarketPriceParams,
) (filterorders_.SortedMarketOrders, error) {
	if rep, err := mpc.structureClient.Fetch(
		ctx,
		filterorders_.FilterStructureMarketParams{
			WebRefreshToken: *params.PricingInfo.MarketRefreshToken,
			StructureId:     params.PricingInfo.MarketLocationId,
			TypeId:          params.TypeId,
			IsBuy:           params.PricingInfo.IsBuy,
		},
	); err != nil {
		return filterorders_.SortedMarketOrders{}, err
	} else {
		return rep.Data(), nil
	}
}

func (mpc MarketPriceClient) fetchStation(
	ctx context.Context,
	params MarketPriceParams,
) (filterorders_.SortedMarketOrders, error) {
	regionId, _ := params.PricingInfo.RegionId()
	if rep, err := mpc.regionClient.Fetch(
		ctx,
		filterorders_.FilterRegionMarketParams{
			RegionId:   regionId,
			TypeId:     params.TypeId,
			IsBuy:      params.PricingInfo.IsBuy,
			LocationId: params.PricingInfo.MarketLocationId,
		},
	); err != nil {
		return filterorders_.SortedMarketOrders{}, err
	} else {
		return rep.Data(), nil
	}
}
