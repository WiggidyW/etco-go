package etcogobucket

import (
	"context"
)

func (bc *BucketClient) ReadCoreData(
	ctx context.Context,
	capacityBuybackSystemTypeMaps int,
	capacityShopLocationTypeMaps int,
	capacityBuybackSystems int,
	capacityShopLocations int,
	capacityBannedFlagSets int,
	capacityPricings int,
	capacityMarkets int,
) (v CoreBucketData, err error) {
	v = CoreBucketData{
		BuybackSystemTypeMaps: make(
			[]BuybackSystemTypeMap,
			0,
			capacityBuybackSystemTypeMaps,
		),
		ShopLocationTypeMaps: make(
			[]ShopLocationTypeMap,
			0,
			capacityShopLocationTypeMaps,
		),
		BuybackSystems: make(
			map[SystemId]BuybackSystem,
			capacityBuybackSystems,
		),
		ShopLocations: make(
			map[LocationId]ShopLocation,
			capacityShopLocations,
		),
		BannedFlagSets: make(
			[]BannedFlagSet,
			0,
			capacityBannedFlagSets,
		),
		Pricings: make(
			[]Pricing,
			0,
			capacityPricings,
		),
		Markets: make(
			[]Market,
			0,
			capacityMarkets,
		),
	}
	_, err = read(
		bc,
		ctx,
		BUILD,
		OBJNAME_CORE_DATA,
		&v,
	)
	return v, err
}

func (bc *BucketClient) WriteCoreData(
	ctx context.Context,
	v CoreBucketData,
) error {
	return write(
		bc,
		ctx,
		BUILD,
		OBJNAME_CORE_DATA,
		v,
	)
}
