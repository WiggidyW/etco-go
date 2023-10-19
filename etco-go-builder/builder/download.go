package builder

import (
	"context"
	"time"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	UPDATER_BUCKET_DATA_TIMEOUT   = 300 * time.Second
	SDE_BUCKET_DATA_TIMEOUT       = 300 * time.Second
	CORE_BUCKET_DATA_TIMEOUT      = 300 * time.Second
	CONSTANTS_BUCKET_DATA_TIMEOUT = 300 * time.Second
)

func transceiveDownloadConstantsBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
	chnSend chanresult.ChanSendResult[b.ConstantsData],
) error {
	constantsBucketData, err := downloadConstantsBucketData(
		ctx,
		bucketClient,
	)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(constantsBucketData)
	}
}
func downloadConstantsBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
) (b.ConstantsData, error) {
	ctx, cancel := context.WithTimeout(ctx, CONSTANTS_BUCKET_DATA_TIMEOUT)
	defer cancel()
	return bucketClient.ReadConstantsData(ctx)
}

// func transceiveDownloadUpdaterBucketData(
//
//	ctx context.Context,
//	bucketClient *b.BucketClient,
//	chnSend chanresult.ChanSendResult[b.UpdaterData],
//
//	) error {
//		updaterBucketData, err := downloadUpdaterBucketData(ctx, bucketClient)
//		if err != nil {
//			return chnSend.SendErr(err)
//		} else {
//			return chnSend.SendOk(updaterBucketData)
//		}
//	}
func downloadUpdaterBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
) (b.UpdaterData, error) {
	ctx, cancel := context.WithTimeout(ctx, UPDATER_BUCKET_DATA_TIMEOUT)
	defer cancel()
	return bucketClient.ReadUpdaterData(ctx)
}

// func transceiveDownloadSDEBucketData(
//
//	ctx context.Context,
//	bucketClient *b.BucketClient,
//	chnSend chanresult.ChanSendResult[b.SDEBucketData],
//
//	) error {
//		sdeBucketData, err := downloadSDEBucketData(ctx, bucketClient)
//		if err != nil {
//			return chnSend.SendErr(err)
//		} else {
//			return chnSend.SendOk(sdeBucketData)
//		}
//	}
func downloadSDEBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
) (b.SDEBucketData, error) {
	ctx, cancel := context.WithTimeout(ctx, SDE_BUCKET_DATA_TIMEOUT)
	defer cancel()
	return bucketClient.ReadSDEData(ctx, 0, 0, 0, 0, 0, 0, 0, 0, 0)
}

// func transceiveDownloadCoreBucketData(
//
//	ctx context.Context,
//	bucketClient *b.BucketClient,
//	chnSend chanresult.ChanSendResult[b.CoreBucketData],
//
//	) error {
//		coreBucketData, err := downloadCoreBucketData(ctx, bucketClient)
//		if err != nil {
//			return chnSend.SendErr(err)
//		} else {
//			return chnSend.SendOk(coreBucketData)
//		}
//	}
func downloadCoreBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
) (b.CoreBucketData, error) {
	ctx, cancel := context.WithTimeout(ctx, CORE_BUCKET_DATA_TIMEOUT)
	defer cancel()
	return bucketClient.ReadCoreData(ctx, 0, 0, 0, 0, 0, 0, 0)
}
