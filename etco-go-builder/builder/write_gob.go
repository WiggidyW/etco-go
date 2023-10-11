package builder

import (
	"context"
	"encoding/gob"
	"os"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"
)

// func transceiveWriteCoreBucketData(
// 	ctx context.Context,
// 	fileDir string,
// 	coreBucketData b.CoreBucketData,
// 	chnSend chanresult.ChanSendResult[struct{}],
// ) error {
// 	err := writeCoreBucketData(ctx, fileDir, coreBucketData)
// 	if err != nil {
// 		return chnSend.SendErr(err)
// 	} else {
// 		return chnSend.SendOk(struct{}{})
// 	}
// }

func writeCoreBucketData(
	ctx context.Context,
	fileDir string,
	coreBucketData b.CoreBucketData,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnSend, chnRecv := chanresult.
		NewChanResult[struct{}](ctx, 7, 0).Split()
	go transceiveWriteGobFile(
		coreBucketData.BuybackSystemTypeMaps,
		fileDir+"/"+b.FILENAME_CORE_BUYBACK_SYSTEM_TYPE_MAPS,
		chnSend,
	)
	go transceiveWriteGobFile(
		coreBucketData.ShopLocationTypeMaps,
		fileDir+"/"+b.FILENAME_CORE_SHOP_LOCATION_TYPE_MAPS,
		chnSend,
	)
	go transceiveWriteGobFile(
		coreBucketData.BuybackSystems,
		fileDir+"/"+b.FILENAME_CORE_BUYBACK_SYSTEMS,
		chnSend,
	)
	go transceiveWriteGobFile(
		coreBucketData.ShopLocations,
		fileDir+"/"+b.FILENAME_CORE_SHOP_LOCATIONS,
		chnSend,
	)
	go transceiveWriteGobFile(
		coreBucketData.BannedFlagSets,
		fileDir+"/"+b.FILENAME_CORE_BANNED_FLAG_SETS,
		chnSend,
	)
	go transceiveWriteGobFile(
		coreBucketData.Pricings,
		fileDir+"/"+b.FILENAME_CORE_PRICINGS,
		chnSend,
	)
	go transceiveWriteGobFile(
		coreBucketData.Markets,
		fileDir+"/"+b.FILENAME_CORE_MARKETS,
		chnSend,
	)

	for i := 0; i < 7; i++ {
		_, err := chnRecv.Recv()
		if err != nil {
			return err
		}
	}

	return nil
}

// func transceiveWriteSDEBucketData(
// 	ctx context.Context,
// 	fileDir string,
// 	sdeBucketData b.SDEBucketData,
// 	chnSend chanresult.ChanSendResult[struct{}],
// ) error {
// 	err := writeSDEBucketData(ctx, fileDir, sdeBucketData)
// 	if err != nil {
// 		return chnSend.SendErr(err)
// 	} else {
// 		return chnSend.SendOk(struct{}{})
// 	}
// }

func writeSDEBucketData(
	ctx context.Context,
	fileDir string,
	sdeBucketData b.SDEBucketData,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnSend, chnRecv := chanresult.
		NewChanResult[struct{}](ctx, 9, 0).Split()
	go transceiveWriteGobFile(
		sdeBucketData.NameToTypeId,
		fileDir+"/"+b.FILENAME_SDE_NAME_TO_TYPE_ID,
		chnSend,
	)
	go transceiveWriteGobFile(
		sdeBucketData.TypeDataMap,
		fileDir+"/"+b.FILENAME_SDE_TYPE_DATA_MAP,
		chnSend,
	)
	go transceiveWriteGobFile(
		sdeBucketData.Categories,
		fileDir+"/"+b.FILENAME_SDE_CATEGORIES,
		chnSend,
	)
	go transceiveWriteGobFile(
		sdeBucketData.Groups,
		fileDir+"/"+b.FILENAME_SDE_GROUPS,
		chnSend,
	)
	go transceiveWriteGobFile(
		sdeBucketData.MarketGroups,
		fileDir+"/"+b.FILENAME_SDE_MARKET_GROUPS,
		chnSend,
	)
	go transceiveWriteGobFile(
		sdeBucketData.TypeVolumes,
		fileDir+"/"+b.FILENAME_SDE_TYPE_VOLUMES,
		chnSend,
	)
	go transceiveWriteGobFile(
		sdeBucketData.Regions,
		fileDir+"/"+b.FILENAME_SDE_REGIONS,
		chnSend,
	)
	go transceiveWriteGobFile(
		sdeBucketData.Systems,
		fileDir+"/"+b.FILENAME_SDE_SYSTEMS,
		chnSend,
	)
	go transceiveWriteGobFile(
		sdeBucketData.Stations,
		fileDir+"/"+b.FILENAME_SDE_STATIONS,
		chnSend,
	)

	for i := 0; i < 9; i++ {
		_, err := chnRecv.Recv()
		if err != nil {
			return err
		}
	}

	return nil
}

func transceiveWriteGobFile[T any](
	t T,
	filePath string,
	chnSend chanresult.ChanSendResult[struct{}],
) error {
	err := writeGobFile(t, filePath)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(struct{}{})
	}
}

func writeGobFile[T any](
	t T,
	filePath string,
) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(t)
}
