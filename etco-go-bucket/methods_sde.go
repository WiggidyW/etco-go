package etcogobucket

import (
	"context"
)

func (bc *BucketClient) ReadSDEData(
	ctx context.Context,
	capacityCategories int,
	capacityGroups int,
	capacityMarketGroups int,
	capacityTypeVolumes int,
	capacityNameToTypeId int,
	capacityRegions int,
	capacitySystems int,
	capacityStations int,
	capacityTypeDataMap int,
) (v SDEBucketData, err error) {
	v = SDEBucketData{
		Categories: make(
			[]CategoryName,
			0,
			capacityCategories,
		),
		Groups: make(
			[]Group,
			0,
			capacityGroups,
		),
		MarketGroups: make(
			[]MarketGroup,
			0,
			capacityMarketGroups,
		),
		TypeVolumes: make(
			[]TypeVolume,
			0,
			capacityTypeVolumes,
		),
		NameToTypeId: make(
			map[TypeName]TypeId,
			capacityNameToTypeId,
		),
		Regions: make(
			map[RegionId]RegionName,
			capacityRegions,
		),
		Systems: make(
			map[SystemId]System,
			capacitySystems,
		),
		Stations: make(
			map[StationId]Station,
			capacityStations,
		),
		TypeDataMap: make(
			map[TypeId]TypeData,
			capacityTypeDataMap,
		),
		UpdaterData: SDEUpdaterData{},
	}
	_, err = read(
		bc,
		ctx,
		BUILD,
		OBJNAME_SDE_DATA,
		&v,
	)
	return v, err
}

func (bc *BucketClient) WriteSDEData(
	ctx context.Context,
	v SDEBucketData,
) error {
	return write(
		bc,
		ctx,
		BUILD,
		OBJNAME_SDE_DATA,
		v,
	)
}
