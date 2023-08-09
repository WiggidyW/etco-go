package orders

import (
	"context"
	"fmt"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client/modelclient"
	"github.com/WiggidyW/weve-esi/staticdb/sde"
	"github.com/WiggidyW/weve-esi/staticdb/tc"
)

type MarketOrdersClientFetchParams struct {
	PricingInfo *tc.PricingInfo
	TypeId      int32
}

func NewFetchParams(
	pricingInfo *tc.PricingInfo,
	typeId int32,
) MarketOrdersClientFetchParams {
	return MarketOrdersClientFetchParams{
		PricingInfo: pricingInfo,
		TypeId:      typeId,
	}
}

func (f MarketOrdersClientFetchParams) CacheKey() string {
	return "marketorders-" + f.CacheKeyInner()
}

func (f MarketOrdersClientFetchParams) CacheKeyInner() string {
	return fmt.Sprintf(
		"%d-%t-%d",
		f.PricingInfo.MarketLocationId(),
		f.PricingInfo.IsBuy(),
		f.TypeId,
	)
}

type MarketOrdersClient struct {
	stationClient   *modelclient.ClientOrdersRegion
	structureClient *modelclient.ClientOrdersStructure
}

func (mpc *MarketOrdersClient) Fetch(
	ctx context.Context,
	params MarketOrdersClientFetchParams,
) (*cache.ExpirableData[MarketOrders], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if params.PricingInfo.MarketIsStructure() {
		return mpc.FetchStructure(ctx, params)
	} else {
		return mpc.FetchStation(ctx, params)
	}
}

func (mpc *MarketOrdersClient) FetchStation(
	ctx context.Context,
	params MarketOrdersClientFetchParams,
) (*cache.ExpirableData[MarketOrders], error) {
	// station ID
	locationId := params.PricingInfo.MarketLocationId()
	// is buy order
	isBuy := params.PricingInfo.IsBuy()
	// region ID
	stationInfo, ok := sde.KVReaderStations.Get(int32(locationId))
	if !ok {
		return nil, fmt.Errorf("FetchStation: station not found")
	}
	systemInfo, ok := sde.KVReaderSystems.Get(stationInfo.SystemId)
	if !ok {
		return nil, fmt.Errorf("FetchStation: system not found")
	}

	// fetch the ESI orders
	// likely to be one page due to a restrictive query
	pages, maxExpires, err := mpc.stationClient.FetchAll( // should be one page
		ctx,
		modelclient.NewFetchParamsOrdersRegion(
			systemInfo.RegionId,
			&params.TypeId,
			&isBuy,
		),
	)
	if err != nil {
		return nil, err
	}

	// count the number of ESI orders
	var numOrders int = 0
	for _, page := range pages {
		numOrders += len(page.Data())
	}

	// initialize orders
	orders := make([]MarketOrder, 0, numOrders)

	// handle the pages
	for _, page := range pages {
		// update maxExpires
		if page.Expires().After(maxExpires) {
			maxExpires = page.Expires()
		}
		// validate, convert, and append the orders
		for _, esiOrder := range page.Data() {
			if esiOrder.LocationId == locationId {
				orders = append(orders, MarketOrder{
					Price:    esiOrder.Price,
					Quantity: int64(esiOrder.VolumeRemain),
				})
			}
		}
	}

	// return the orders
	data := cache.NewExpirableData(NewMarketOrders(orders), maxExpires)
	return &data, nil
}

func (mpc *MarketOrdersClient) FetchStructure(
	ctx context.Context,
	params MarketOrdersClientFetchParams,
) (*cache.ExpirableData[MarketOrders], error) {
	// structure ID
	locationId := params.PricingInfo.MarketLocationId()
	// refresh token
	refreshToken, ok := params.PricingInfo.MarketRefreshToken()
	if !ok {
		return nil, fmt.Errorf("FetchStructure: no refresh token")
	}
	// is buy order
	isBuy := params.PricingInfo.IsBuy()

	// fetch the ESI orders
	strm := mpc.structureClient.FetchStreamBlocking( // should be one page
		ctx,
		modelclient.NewFetchParamsOrdersStructure(
			locationId,
			refreshToken,
		),
	)
	defer strm.Close()

	// initialize orders
	orders := make([]MarketOrder, 0, mpc.stationClient.EntriesPerPage())

	// handle the pages
	var maxExpires time.Time = strm.HeadExpires() // head expiry
	for pages := strm.NumPages(); pages > 0; pages-- {
		page, err := strm.Recv()
		if err != nil {
			return nil, err
		}
		// update maxExpires
		if page.Expires().After(maxExpires) {
			maxExpires = page.Expires()
		}
		// validate, convert, and append the orders
		for _, esiOrder := range page.Data() {
			if esiOrder.IsBuyOrder == isBuy &&
				esiOrder.TypeId == params.TypeId {
				orders = append(orders, MarketOrder{
					Price:    esiOrder.Price,
					Quantity: int64(esiOrder.VolumeRemain),
				})
			}
		}
	}

	// return the orders and their expiration date
	data := cache.NewExpirableData(NewMarketOrders(orders), maxExpires)
	return &data, nil
}
