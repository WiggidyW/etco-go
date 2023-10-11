package updater

import (
	"context"
	"sync"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	updaterwtc "github.com/WiggidyW/etco-go-updater/webtocore"
)

func UpdateCoreIfModified(
	ctx context.Context,
	bucketClient *b.BucketClient,
	syncUpdaterData *SyncUpdaterData,
	webAttrs updaterwtc.WebAttrs,
	chnSendDone chanresult.ChanSendResult[struct{}],
) (updating bool) {
	syncUpdaterData.RLock()
	updating = webAttrs.CHECKSUM_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER !=
		syncUpdaterData.CHECKSUM_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER ||
		webAttrs.CHECKSUM_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER !=
			syncUpdaterData.CHECKSUM_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER ||
		webAttrs.CHECKSUM_WEB_BUYBACK_SYSTEMS !=
			syncUpdaterData.CHECKSUM_WEB_BUYBACK_SYSTEMS ||
		webAttrs.CHECKSUM_WEB_SHOP_LOCATIONS !=
			syncUpdaterData.CHECKSUM_WEB_SHOP_LOCATIONS ||
		webAttrs.CHECKSUM_WEB_MARKETS !=
			syncUpdaterData.CHECKSUM_WEB_MARKETS
	syncUpdaterData.RUnlock()

	if !updating {
		return false
	}

	go TransceiveUpdateCore(
		ctx,
		bucketClient,
		syncUpdaterData,
		webAttrs,
		chnSendDone,
	)
	return true
}

func TransceiveUpdateCore(
	ctx context.Context,
	bucketClient *b.BucketClient,
	syncUpdaterData *SyncUpdaterData,
	webAttrs updaterwtc.WebAttrs,
	chnSendDone chanresult.ChanSendResult[struct{}],
) error {
	err := UpdateCore(
		ctx,
		bucketClient,
		syncUpdaterData,
		webAttrs,
	)
	if err != nil {
		return chnSendDone.SendErr(err)
	} else {
		return chnSendDone.SendOk(struct{}{})
	}
}

func UpdateCore(
	ctx context.Context,
	bucketClient *b.BucketClient,
	syncUpdaterData *SyncUpdaterData,
	webAttrs updaterwtc.WebAttrs,
) error {
	wtcBucketData, err := updaterwtc.DownloadAndConvert(
		ctx,
		bucketClient,
	)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		updateWTCSyncUpdaterData(
			syncUpdaterData,
			webAttrs,
			wtcBucketData,
		)
		wg.Done()
	}()

	// avoid writing anything if an error occured somewhere
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// write the data to the bucket
	err = bucketClient.WriteCoreData(ctx, wtcBucketData.CoreBucketData)
	if err != nil {
		return err
	}

	wg.Wait()
	return nil
}

func updateWTCSyncUpdaterData(
	syncUpdaterData *SyncUpdaterData,
	webAttrs updaterwtc.WebAttrs,
	wtcBucketData updaterwtc.WebAndCoreBucketData,
) {
	syncUpdaterData.WLock()
	defer syncUpdaterData.WUnlock()

	syncUpdaterData.CHECKSUM_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER =
		webAttrs.CHECKSUM_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER
	syncUpdaterData.CHECKSUM_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER =
		webAttrs.CHECKSUM_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER
	syncUpdaterData.CHECKSUM_WEB_BUYBACK_SYSTEMS =
		webAttrs.CHECKSUM_WEB_BUYBACK_SYSTEMS
	syncUpdaterData.CHECKSUM_WEB_SHOP_LOCATIONS =
		webAttrs.CHECKSUM_WEB_SHOP_LOCATIONS
	syncUpdaterData.CHECKSUM_WEB_MARKETS =
		webAttrs.CHECKSUM_WEB_MARKETS

	syncUpdaterData.CAPACITY_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER =
		len(wtcBucketData.WebBucketData.BuybackSystemTypeMapsBuilder)
	syncUpdaterData.CAPACITY_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER =
		len(wtcBucketData.WebBucketData.ShopLocationTypeMapsBuilder)
	syncUpdaterData.CAPACITY_WEB_BUYBACK_SYSTEMS =
		len(wtcBucketData.WebBucketData.BuybackSystems)
	syncUpdaterData.CAPACITY_WEB_SHOP_LOCATIONS =
		len(wtcBucketData.WebBucketData.ShopLocations)
	syncUpdaterData.CAPACITY_WEB_MARKETS =
		len(wtcBucketData.WebBucketData.Markets)

	syncUpdaterData.CAPACITY_CORE_BUYBACK_SYSTEM_TYPE_MAPS =
		len(wtcBucketData.CoreBucketData.BuybackSystemTypeMaps)
	syncUpdaterData.CAPACITY_CORE_SHOP_LOCATION_TYPE_MAPS =
		len(wtcBucketData.CoreBucketData.ShopLocationTypeMaps)
	syncUpdaterData.CAPACITY_CORE_BUYBACK_SYSTEMS =
		len(wtcBucketData.CoreBucketData.BuybackSystems)
	syncUpdaterData.CAPACITY_CORE_SHOP_LOCATIONS =
		len(wtcBucketData.CoreBucketData.ShopLocations)
	syncUpdaterData.CAPACITY_CORE_BANNED_FLAG_SETS =
		len(wtcBucketData.CoreBucketData.BannedFlagSets)
	syncUpdaterData.CAPACITY_CORE_PRICINGS =
		len(wtcBucketData.CoreBucketData.Pricings)
	syncUpdaterData.CAPACITY_CORE_MARKETS =
		len(wtcBucketData.CoreBucketData.Markets)

	syncUpdaterData.VERSION_BUYBACK = webAttrs.VERSION_STRING
	syncUpdaterData.VERSION_SHOP = webAttrs.VERSION_STRING
}
