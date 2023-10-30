package builder

import (
	"context"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"
)

func transceiveDownloadAndWriteSDEBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
	gobFileDir string,
	chnSendUpdaterData chanresult.ChanSendResult[b.SDEUpdaterData],
	chnSendDone chanresult.ChanSendResult[struct{}],
) error {
	err := downloadAndWriteSDEBucketData(
		ctx,
		bucketClient,
		gobFileDir,
		chnSendUpdaterData,
	)
	if err != nil {
		return chnSendDone.SendErr(err)
	} else {
		return chnSendDone.SendOk(struct{}{})
	}
}
func downloadAndWriteSDEBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
	gobFileDir string,
	chnSendUpdaterData chanresult.ChanSendResult[b.SDEUpdaterData],
) (err error) {
	sdeBucketData, err := downloadSDEBucketData(ctx, bucketClient)
	if err != nil {
		go chnSendUpdaterData.SendErr(err)
		return err
	} else {
		go chnSendUpdaterData.SendOk(sdeBucketData.UpdaterData)
		return writeSDEBucketData(ctx, gobFileDir, sdeBucketData)
	}
}

func transceiveDownloadAndWriteCoreBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
	gobFileDir string,
	chnSendUpdaterData chanresult.ChanSendResult[b.CoreUpdaterData],
	chnSendDone chanresult.ChanSendResult[struct{}],
) error {
	err := downloadAndWriteCoreBucketData(
		ctx,
		bucketClient,
		gobFileDir,
		chnSendUpdaterData,
	)
	if err != nil {
		return chnSendDone.SendErr(err)
	} else {
		return chnSendDone.SendOk(struct{}{})
	}
}
func downloadAndWriteCoreBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
	gobFileDir string,
	chnSendUpdaterData chanresult.ChanSendResult[b.CoreUpdaterData],
) (err error) {
	coreBucketData, err := downloadCoreBucketData(ctx, bucketClient)
	if err != nil {
		go chnSendUpdaterData.SendErr(err)
		return err
	} else {
		go chnSendUpdaterData.SendOk(coreBucketData.UpdaterData)
		return writeCoreBucketData(ctx, gobFileDir, coreBucketData)
	}
}
