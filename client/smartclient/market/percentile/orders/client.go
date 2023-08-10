package orders

import (
	"context"
	"fmt"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client/modelclient"
	"github.com/WiggidyW/weve-esi/staticdb"
	"github.com/WiggidyW/weve-esi/staticdb/inner/sde"
)

type MrktOrdersParams struct {
	PricingInfo staticdb.PricingInfo
	TypeId      int32
}

func (f MrktOrdersParams) CacheKey() string {
	return "mrktorders-" + f.CacheKeyInner()
}

func (f MrktOrdersParams) CacheKeyInner() string {
	return fmt.Sprintf(
		"%d-%d-%t",
		f.PricingInfo.MrktLocationId,
		f.TypeId,
		f.PricingInfo.IsBuy,
	)
}

type MrktOrdersClient struct {
	stationClient   *modelclient.ClientOrdersRegion
	structureClient *modelclient.ClientOrdersStructure
}

func (mpc *MrktOrdersClient) Fetch(
	ctx context.Context,
	params MrktOrdersParams,
) (*cache.ExpirableData[MrktOrders], error) {
	if params.PricingInfo.MrktIsStructure {
		return mpc.FetchStructure(ctx, params)
	} else {
		return mpc.FetchStation(ctx, params)
	}
}

func (mpc *MrktOrdersClient) FetchStation(
	ctx context.Context,
	params MrktOrdersParams,
) (*cache.ExpirableData[MrktOrders], error) {
	// get the region ID
	var regionID int32
	if station, ok := sde.KVReaderStations.Get(
		int32(params.PricingInfo.MrktLocationId),
	); ok {
		if system, ok := sde.KVReaderSystems.Get(station.SystemId); ok {
			regionID = system.RegionId
		} else {
			return nil, fmt.Errorf("FetchStation: system not found")
		}
	} else {
		return nil, fmt.Errorf("FetchStation: station not found")
	}

	// // fetch the ESI orders

	// likely to be one page due to a restrictive query
	pages, maxExpires, err := mpc.stationClient.FetchAll( // should be one page
		ctx,
		modelclient.NewFetchParamsOrdersRegion(
			regionID,
			&params.TypeId,
			&params.PricingInfo.IsBuy,
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

	// collect the orders from the pages
	orders := make([]MrktOrder, 0, numOrders)
	for _, page := range pages {

		// update maxExpires
		if page.Expires().After(maxExpires) {
			maxExpires = page.Expires()
		}

		// validate, convert, and append the orders
		for _, esiOrder := range page.Data() {
			if esiOrder.LocationId == params.PricingInfo.
				MrktLocationId {
				orders = append(orders, MrktOrder{
					Price:    esiOrder.Price,
					Quantity: int64(esiOrder.VolumeRemain),
				})
			}
		}
	}

	// return the orders
	data := cache.NewExpirableData(NewMrktOrders(orders), maxExpires)
	return &data, nil
}

func (mpc *MrktOrdersClient) FetchStructure(
	ctx context.Context,
	params MrktOrdersParams,
) (*cache.ExpirableData[MrktOrders], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the ESI orders
	chn, err := mpc.structureClient.FetchStreamBlocking(
		ctx,
		modelclient.NewFetchParamsOrdersStructure(
			params.PricingInfo.MrktLocationId,
			*params.PricingInfo.MrktRefreshToken,
		),
	)
	if err != nil {
		return nil, err
	}

	// // collect the orders
	orders := make([]MrktOrder, 0, mpc.stationClient.EntriesPerPage())
	var maxExpires time.Time = chn.HeadExpires() // initialize as head expires

	for pages := chn.NumPages(); pages > 0; pages-- {
		page, err := chn.Recv()
		if err != nil {
			return nil, err
		}

		// update maxExpires
		if page.Expires().After(maxExpires) {
			maxExpires = page.Expires()
		}

		// validate, convert, and append the orders
		for _, esiOrder := range page.Data() {
			if esiOrder.IsBuyOrder == params.PricingInfo.IsBuy &&
				esiOrder.TypeId == params.TypeId {
				orders = append(orders, MrktOrder{
					Price:    esiOrder.Price,
					Quantity: int64(esiOrder.VolumeRemain),
				})
			}
		}
	}

	// return the orders and their expiration date
	data := cache.NewExpirableData(NewMrktOrders(orders), maxExpires)
	return &data, nil
}
