package builder

import (
	"context"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"
)

func transceiveDownloadAndWriteUpdaterAndConstantsBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
	constantsFilePath string,
	chnSend chanresult.ChanSendResult[struct{}],
) error {
	err := downloadAndWriteUpdaterAndConstantsBucketData(
		ctx,
		bucketClient,
		constantsFilePath,
	)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(struct{}{})
	}
}

func downloadAndWriteUpdaterAndConstantsBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
	constantsFilePath string,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnSend, chnRecv := chanresult.
		NewChanResult[b.ConstantsData](ctx, 2, 0).Split()
	go transceiveDownloadConstantsBucketData(ctx, bucketClient, chnSend)

	updaterBucketData, err := downloadUpdaterBucketData(ctx, bucketClient)
	if err != nil {
		return err
	}
	constantsBucketData, err := chnRecv.Recv()
	if err != nil {
		return err
	}

	return writeConstants(
		constantsFilePath,
		constantsBucketData,
		updaterBucketData,
	)
}

func transceiveDownloadAndWriteSDEBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
	gobFileDir string,
	chnSend chanresult.ChanSendResult[struct{}],
) error {
	err := downloadAndWriteSDEBucketData(ctx, bucketClient, gobFileDir)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(struct{}{})
	}
}

func downloadAndWriteSDEBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
	gobFileDir string,
) error {
	sdeBucketData, err := downloadSDEBucketData(ctx, bucketClient)
	if err != nil {
		return err
	}
	return writeSDEBucketData(ctx, gobFileDir, sdeBucketData)
}

func transceiveDownloadAndWriteCoreBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
	gobFileDir string,
	chnSend chanresult.ChanSendResult[struct{}],
) error {
	err := downloadAndWriteCoreBucketData(ctx, bucketClient, gobFileDir)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(struct{}{})
	}
}

func downloadAndWriteCoreBucketData(
	ctx context.Context,
	bucketClient *b.BucketClient,
	gobFileDir string,
) error {
	coreBucketData, err := downloadCoreBucketData(ctx, bucketClient)
	if err != nil {
		return err
	}
	return writeCoreBucketData(ctx, gobFileDir, coreBucketData)
}
