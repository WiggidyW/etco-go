package primarysdedata

import (
	"context"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"
)

func convert(
	ctx context.Context,
	rawSDEData RawSDEData,
) (
	etcoPrimarySDEData PrimarySDEData,
	err error,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnStations := chanresult.
		NewChanResult[map[b.StationId]b.Station](ctx, 1, 0)
	go transceiveConvertSDEStations(
		rawSDEData.BSDStaStations,
		chnStations.ToSend(),
	)

	feTypeDatas, err := filterExtendSDETypeDataMap(rawSDEData.FSDTypeIds)
	if err != nil {
		return etcoPrimarySDEData, err
	}

	chnNameToTypeId := chanresult.
		NewChanResult[map[string]b.TypeId](ctx, 1, 0)
	go transceiveConvertFETypeDatasToETCONameToTypeId(
		feTypeDatas,
		chnNameToTypeId.ToSend(),
	)

	etcoCategories := make([]b.CategoryName, 0)
	etcoCategoriesPtr := &etcoCategories
	etcoCategoriesindexMap := make(map[CategoryId]int)

	etcoGroups := make([]b.Group, 0)
	etcoGroupsPtr := &etcoGroups
	etcoGroupsIndexMap := make(map[GroupId]int)

	etcoMarketGroups := make([]b.MarketGroup, 0)
	etcoMarketGroupsPtr := &etcoMarketGroups
	etcoMarketGroupsIndexMap := make(map[MarketGroupId]int)

	etcoTypeVolumes := make([]b.TypeVolume, 0)
	etcoTypeVolumesPtr := &etcoTypeVolumes
	etcoTypeVolumesIndexMap := make(map[float64]int)

	for _, feTypeData := range feTypeDatas {

		if err = feTypeData.addSDEGroup(
			rawSDEData.FSDGroupIds,
			rawSDEData.FSDCategoryIds,
			etcoGroupsPtr,
			etcoGroupsIndexMap,
			etcoCategoriesPtr,
			etcoCategoriesindexMap,
		); err != nil {
			return etcoPrimarySDEData, err
		}

		if err = feTypeData.addSDEMarketGroup(
			rawSDEData.FSDMarketGroups,
			etcoMarketGroupsPtr,
			etcoMarketGroupsIndexMap,
		); err != nil {
			return etcoPrimarySDEData, err
		}

		if err = feTypeData.addSDETypeMaterials(
			rawSDEData.FSDTypeMaterials,
		); err != nil {
			return etcoPrimarySDEData, err
		}

		feTypeData.addSDETypeVolume(
			etcoTypeVolumesPtr,
			etcoTypeVolumesIndexMap,
		)
	}

	etcoTypeDataMap, err := convertFETypeDatasToETCOTypeDataMap(
		feTypeDatas,
	)
	if err != nil {
		return etcoPrimarySDEData, err
	}

	etcoNameToTypeId, err := chnNameToTypeId.Recv()
	if err != nil {
		return etcoPrimarySDEData, err
	}

	etcoStations, err := chnStations.Recv()
	if err != nil {
		return etcoPrimarySDEData, err
	}

	return PrimarySDEData{
		ETCOCategories:   etcoCategories,
		ETCOGroups:       etcoGroups,
		ETCOMarketGroups: etcoMarketGroups,
		ETCOTypeVolumes:  etcoTypeVolumes,
		ETCONameToTypeId: etcoNameToTypeId,
		ETCOStations:     etcoStations,
		ETCOTypeDataMap:  etcoTypeDataMap,
	}, nil
}
