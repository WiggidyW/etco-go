package primarysdedata

import (
	"context"
	"fmt"

	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go-updater/sde/loadyaml"
)

const (
	AFFIX_BSD_STASTATIONS   = "bsd/staStations.yaml"
	AFFIX_FSD_CATEGORYIDS   = "fsd/categoryIDs.yaml"
	AFFIX_FSD_TYPEMATERIALS = "fsd/typeMaterials.yaml"
	AFFIX_FSD_TYPEIDS       = "fsd/typeIDs.yaml"
	AFFIX_FSD_MARKETGROUPS  = "fsd/marketGroups.yaml"
	AFFIX_FSD_GROUPIDS      = "fsd/groupIDs.yaml"
)

func loadRawSDEData(
	ctx context.Context,
	pathSDE string,
) (RawSDEData, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnFSDTypeIDs := chanresult.NewChanResult[FSDTypeIds](ctx, 1, 0)
	go transceiveLoadRawSDEData(
		chnFSDTypeIDs.ToSend(),
		fmt.Sprintf("%s/%s", pathSDE, AFFIX_FSD_TYPEIDS),
	)

	chnFSDTypeMaterials := chanresult.
		NewChanResult[FSDTypeMaterials](ctx, 1, 0)
	go transceiveLoadRawSDEData(
		chnFSDTypeMaterials.ToSend(),
		fmt.Sprintf("%s/%s", pathSDE, AFFIX_FSD_TYPEMATERIALS),
	)

	chnFSDMarketGroups := chanresult.
		NewChanResult[FSDMarketGroups](ctx, 1, 0)
	go transceiveLoadRawSDEData(
		chnFSDMarketGroups.ToSend(),
		fmt.Sprintf("%s/%s", pathSDE, AFFIX_FSD_MARKETGROUPS),
	)

	chnFSDGroupIDs := chanresult.NewChanResult[FSDGroupIds](ctx, 1, 0)
	go transceiveLoadRawSDEData(
		chnFSDGroupIDs.ToSend(),
		fmt.Sprintf("%s/%s", pathSDE, AFFIX_FSD_GROUPIDS),
	)

	chnFSDCategoryIDs := chanresult.
		NewChanResult[FSDCategoryIds](ctx, 1, 0)
	go transceiveLoadRawSDEData(
		chnFSDCategoryIDs.ToSend(),
		fmt.Sprintf("%s/%s", pathSDE, AFFIX_FSD_CATEGORYIDS),
	)

	chnBSDStaStations := chanresult.
		NewChanResult[BSDStaStations](ctx, 1, 0)
	go transceiveLoadRawSDEData(
		chnBSDStaStations.ToSend(),
		fmt.Sprintf("%s/%s", pathSDE, AFFIX_BSD_STASTATIONS),
	)

	if rawSdeData, err := chanresult.RecvOneOfEach6(
		chnFSDTypeIDs.ToRecv(),
		chnFSDTypeMaterials.ToRecv(),
		chnFSDMarketGroups.ToRecv(),
		chnFSDGroupIDs.ToRecv(),
		chnFSDCategoryIDs.ToRecv(),
		chnBSDStaStations.ToRecv(),
	); err != nil {
		return RawSDEData{}, err
	} else {
		return RawSDEData{
			FSDTypeIds:       rawSdeData.T1,
			FSDTypeMaterials: rawSdeData.T2,
			FSDMarketGroups:  rawSdeData.T3,
			FSDGroupIds:      rawSdeData.T4,
			FSDCategoryIds:   rawSdeData.T5,
			BSDStaStations:   rawSdeData.T6,
		}, nil
	}
}

func transceiveLoadRawSDEData[RD any](
	chnSendRawSDEData chanresult.ChanSendResult[RD],
	path string,
) error {
	if rawSDEData, err := loadyaml.LoadYaml[RD](path); err != nil {
		return chnSendRawSDEData.SendErr(err)
	} else {
		return chnSendRawSDEData.SendOk(rawSDEData)
	}
}
