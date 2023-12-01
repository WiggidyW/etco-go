package sde

import (
	"context"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	pd "github.com/WiggidyW/etco-go-updater/sde/primarysdedata"
	ud "github.com/WiggidyW/etco-go-updater/sde/universesdedata"
)

func LoadAndConvert(
	ctx context.Context,
	sdeChecksum string,
	pathSDE string,
) (
	sdeBucketData b.SDEBucketData,
	err error,
) {
	chnPrimaryData := chanresult.
		NewChanResult[pd.PrimarySDEData](ctx, 1, 0)
	go pd.TransceiveLoadAndConvert(
		ctx,
		pathSDE,
		chnPrimaryData.ToSend(),
	)

	chnUniverseData := chanresult.
		NewChanResult[ud.UniverseSDEData](ctx, 1, 0)
	go ud.TransceiveLoadAndConvert(pathSDE, chnUniverseData.ToSend())

	primaryData, err := chnPrimaryData.Recv()
	if err != nil {
		return sdeBucketData, err
	}

	universeData, err := chnUniverseData.Recv()
	if err != nil {
		return sdeBucketData, err
	}

	sdeBucketData = b.SDEBucketData{
		Categories:   primaryData.ETCOCategories,
		Groups:       primaryData.ETCOGroups,
		MarketGroups: primaryData.ETCOMarketGroups,
		TypeVolumes:  primaryData.ETCOTypeVolumes,
		NameToTypeId: primaryData.ETCONameToTypeId,
		Stations:     primaryData.ETCOStations,
		TypeDataMap:  primaryData.ETCOTypeDataMap,
		Regions:      universeData.ETCORegions,
		Systems:      universeData.ETCOSystems,
		SystemIds:    universeData.ETCOSystemIds,
		UpdaterData: b.SDEUpdaterData{
			CHECKSUM_SDE:                 sdeChecksum,
			CAPACITY_SDE_CATEGORIES:      len(primaryData.ETCOCategories),
			CAPACITY_SDE_GROUPS:          len(primaryData.ETCOGroups),
			CAPACITY_SDE_MARKET_GROUPS:   len(primaryData.ETCOMarketGroups),
			CAPACITY_SDE_TYPE_VOLUMES:    len(primaryData.ETCOTypeVolumes),
			CAPACITY_SDE_NAME_TO_TYPE_ID: len(primaryData.ETCONameToTypeId),
			CAPACITY_SDE_STATIONS:        len(primaryData.ETCOStations),
			CAPACITY_SDE_TYPE_DATA_MAP:   len(primaryData.ETCOTypeDataMap),
			CAPACITY_SDE_REGIONS:         len(universeData.ETCORegions),
			CAPACITY_SDE_SYSTEMS:         len(universeData.ETCOSystems),
			CAPACITY_SDE_SYSTEM_IDS:      len(universeData.ETCOSystemIds),
		},
	}
	return sdeBucketData, nil
}
