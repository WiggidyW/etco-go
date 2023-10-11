package sde

import (
	"context"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	pd "github.com/WiggidyW/etco-go-updater/sde/primarysdedata"
	ud "github.com/WiggidyW/etco-go-updater/sde/universesdedata"
)

func TransceiveLoadAndConvert(
	ctx context.Context,
	pathSDE string,
	chnSend chanresult.ChanSendResult[b.SDEBucketData],
) error {
	etcoSDEBucketData, err := LoadAndConvert(ctx, pathSDE)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(etcoSDEBucketData)
	}
}

func LoadAndConvert(
	ctx context.Context,
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

	return b.SDEBucketData{
		Categories:   primaryData.ETCOCategories,
		Groups:       primaryData.ETCOGroups,
		MarketGroups: primaryData.ETCOMarketGroups,
		TypeVolumes:  primaryData.ETCOTypeVolumes,
		NameToTypeId: primaryData.ETCONameToTypeId,
		Stations:     primaryData.ETCOStations,
		TypeDataMap:  primaryData.ETCOTypeDataMap,
		Regions:      universeData.ETCORegions,
		Systems:      universeData.ETCOSystems,
	}, nil
}
