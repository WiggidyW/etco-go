package updater

import (
	"context"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	updaterwtc "github.com/WiggidyW/etco-go-updater/webtocore"
)

func transceiveCoreModified(
	ctx context.Context,
	bucketClient *b.BucketClient,
	skipCore bool,
	chnSendModified chanresult.ChanSendResult[*updaterwtc.WebAttrs],
) error {
	modified, err := CoreModified(ctx, bucketClient, skipCore)
	if err != nil {
		return chnSendModified.SendErr(err)
	} else {
		return chnSendModified.SendOk(modified)
	}
}

func CoreModified(
	ctx context.Context,
	bucketClient *b.BucketClient,
	skipCore bool,
) (
	modified *updaterwtc.WebAttrs, // Some if modified, nil if not
	err error,
) {
	if skipCore {
		return nil, nil
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
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	)
	if err != nil {
		return nil, err
	}

	// receive the web attrs
	webAttrs, err := chnRecvWebAttrs.Recv()
	if err != nil {
		return nil, err
	} else if webAttrs.CHECKSUM_WEB_BUYBACK_SYSTEMS ==
		prevCoreBucketData.UpdaterData.CHECKSUM_WEB_BUYBACK_SYSTEMS &&
		webAttrs.CHECKSUM_WEB_SHOP_LOCATIONS ==
			prevCoreBucketData.UpdaterData.CHECKSUM_WEB_SHOP_LOCATIONS &&
		webAttrs.CHECKSUM_WEB_MARKETS ==
			prevCoreBucketData.UpdaterData.CHECKSUM_WEB_MARKETS &&
		webAttrs.CHECKSUM_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER ==
			prevCoreBucketData.UpdaterData.CHECKSUM_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER &&
		webAttrs.CHECKSUM_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER ==
			prevCoreBucketData.UpdaterData.CHECKSUM_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER &&
		webAttrs.CHECKSUM_WEB_HAUL_ROUTE_TYPE_MAPS_BUILDER ==
			prevCoreBucketData.UpdaterData.CHECKSUM_WEB_HAUL_ROUTE_TYPE_MAPS_BUILDER &&
		webAttrs.CHECKSUM_WEB_HAUL_ROUTES ==
			prevCoreBucketData.UpdaterData.CHECKSUM_WEB_HAUL_ROUTES {
		// return false if all checksums match
		return nil, nil
	}

	return &webAttrs, nil
}

func UpdateCore(
	ctx context.Context,
	bucketClient *b.BucketClient,
	webAttrs updaterwtc.WebAttrs,
	sdeSystems map[b.SystemId]b.System,
) error {
	// download and convert Web bucket data into writeable core bucket data
	coreBucketData, err := updaterwtc.DownloadAndConvert(
		ctx,
		bucketClient,
		webAttrs,
		sdeSystems,
	)
	if err != nil {
		return err
	}

	// write the data to the bucket
	return bucketClient.WriteCoreData(ctx, coreBucketData)
}
