package primarysdedata

import (
	"context"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"
)

type PrimarySDEData struct {
	ETCOCategories   []b.CategoryName
	ETCOGroups       []b.Group
	ETCOMarketGroups []b.MarketGroup
	ETCOTypeVolumes  []b.TypeVolume
	ETCONameToTypeId map[string]b.TypeId
	ETCOStations     map[b.StationId]b.Station
	ETCOTypeDataMap  map[b.TypeId]b.TypeData
}

func TransceiveLoadAndConvert(
	ctx context.Context,
	pathSDE string,
	chnSend chanresult.ChanSendResult[PrimarySDEData],
) error {
	etcoPrimarySDEData, err := LoadAndConvert(ctx, pathSDE)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(etcoPrimarySDEData)
	}
}

func LoadAndConvert(
	ctx context.Context,
	pathSDE string,
) (
	etcoPrimarySDEData PrimarySDEData,
	err error,
) {
	rawSDEData, err := loadRawSDEData(ctx, pathSDE)
	if err != nil {
		return etcoPrimarySDEData, err
	}
	etcoPrimarySDEData, err = convert(ctx, rawSDEData)
	if err != nil {
		return etcoPrimarySDEData, err
	}
	return etcoPrimarySDEData, nil
}
