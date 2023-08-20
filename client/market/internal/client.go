package internal

import (
	"context"
	"fmt"

	"github.com/WiggidyW/weve-esi/client/market/internal/filter"
	reg "github.com/WiggidyW/weve-esi/client/market/internal/filter/region"
	stc "github.com/WiggidyW/weve-esi/client/market/internal/filter/structure"
	"github.com/WiggidyW/weve-esi/logger"
)

type MarketPriceClient struct {
	structureClient stc.WC_FilterStructureMarketClient
	regionClient    reg.WC_FilterRegionMarketClient
}

// returns nil if there are no orders to be found
func (mpc MarketPriceClient) Fetch(
	ctx context.Context,
	params MarketPriceParams,
) (*PositivePrice, error) {
	var sortedOrders filter.SortedMarketOrders
	var err error

	if params.PricingInfo.MrktIsStructure {
		sortedOrders, err = mpc.fetchStructure(ctx, params)
	} else {
		sortedOrders, err = mpc.fetchStation(ctx, params)
	}
	if err != nil {
		return nil, err
	} else if sortedOrders.Quantity == 0 {
		return nil, nil
	}

	price := sortedOrders.PercentilePrice(params.PricingInfo.Prctile)
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
) (filter.SortedMarketOrders, error) {
	if rep, err := mpc.structureClient.Fetch(
		ctx,
		stc.FilterStructureMarketParams{
			WebRefreshToken: *params.PricingInfo.MrktRefreshToken,
			StructureId:     params.PricingInfo.MrktLocationId,
			TypeId:          params.TypeId,
			IsBuy:           params.PricingInfo.IsBuy,
		},
	); err != nil {
		return filter.SortedMarketOrders{}, err
	} else {
		return rep.Data(), nil
	}
}

func (mpc MarketPriceClient) fetchStation(
	ctx context.Context,
	params MarketPriceParams,
) (filter.SortedMarketOrders, error) {
	regionId, _ := params.PricingInfo.RegionId()
	if rep, err := mpc.regionClient.Fetch(
		ctx,
		reg.FilterRegionMarketParams{
			RegionId:   regionId,
			TypeId:     params.TypeId,
			IsBuy:      params.PricingInfo.IsBuy,
			LocationId: params.PricingInfo.MrktLocationId,
		},
	); err != nil {
		return filter.SortedMarketOrders{}, err
	} else {
		return rep.Data(), nil
	}
}
