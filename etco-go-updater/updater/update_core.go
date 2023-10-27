package updater

import (
	"context"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	updaterwtc "github.com/WiggidyW/etco-go-updater/webtocore"
)

func transceiveUpdateCoreIfModified(
	ctx context.Context,
	bucketClient *b.BucketClient,
	skipCore bool,
	chnSendModified chanresult.ChanSendResult[bool],
) error {
	modified, err := UpdateCoreIfModified(
		ctx,
		bucketClient,
		skipCore,
	)
	if err != nil {
		return chnSendModified.SendErr(err)
	} else {
		return chnSendModified.SendOk(modified)
	}
}
func UpdateCoreIfModified(
	ctx context.Context,
	bucketClient *b.BucketClient,
	skipCore bool,
) (
	modified bool,
	err error,
) {
	if skipCore {
		return false, nil
	}

	// fetch the web attrs in a goroutine
	chnSendWebAttrs, chnRecvWebAttrs :=
		chanresult.NewChanResult[updaterwtc.WebAttrs](ctx, 1, 0).Split()
	go updaterwtc.TransceiveDownloadWebAttrs(
		ctx,
		bucketClient,
		chnSendWebAttrs,
	)

	// fetch previous Core data
	prevCoreBucketData, err := bucketClient.ReadCoreData(
		ctx,
		0, 0, 0, 0, 0, 0, 0,
	)
	if err != nil {
		return false, err
	}

	// receive the web attrs
	webAttrs, err := chnRecvWebAttrs.Recv()
	if err != nil {
		return false, err
	} else if webAttrs.CHECKSUM_WEB_BUYBACK_SYSTEMS ==
		prevCoreBucketData.UpdaterData.CHECKSUM_WEB_BUYBACK_SYSTEMS &&
		webAttrs.CHECKSUM_WEB_SHOP_LOCATIONS ==
			prevCoreBucketData.UpdaterData.CHECKSUM_WEB_SHOP_LOCATIONS &&
		webAttrs.CHECKSUM_WEB_MARKETS ==
			prevCoreBucketData.UpdaterData.CHECKSUM_WEB_MARKETS &&
		webAttrs.CHECKSUM_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER ==
			prevCoreBucketData.UpdaterData.CHECKSUM_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER &&
		webAttrs.CHECKSUM_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER ==
			prevCoreBucketData.UpdaterData.CHECKSUM_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER {
		// return false if all checksums match
		return false, nil
	}

	// download, convert, and write Core build bucket data
	err = UpdateCore(
		ctx,
		bucketClient,
		webAttrs,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdateCore(
	ctx context.Context,
	bucketClient *b.BucketClient,
	webAttrs updaterwtc.WebAttrs,
) error {
	// download and convert Web bucket data into writeable core bucket data
	coreBucketData, err := updaterwtc.DownloadAndConvert(
		ctx,
		bucketClient,
		webAttrs,
	)
	if err != nil {
		return err
	}

	// write the data to the bucket
	return bucketClient.WriteCoreData(ctx, coreBucketData)
}
