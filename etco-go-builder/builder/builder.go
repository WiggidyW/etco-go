package builder

import (
	"context"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"
)

func DownloadAndWrite(
	ctx context.Context,
	bucketClient *b.BucketClient,
	gobFileDir string,
	constantsFilePath string,
) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnSendDone, chnRecvDone := chanresult.
		NewChanResult[struct{}](ctx, 3, 0).Split()

	// Download and write SDE .gob files
	// Also, channel-send updater data for use in constants file generation
	chnSendSDEUpdaterData, chnRecvSDEUpdaterData := chanresult.
		NewChanResult[b.SDEUpdaterData](ctx, 1, 0).Split()
	go transceiveDownloadAndWriteSDEBucketData(
		ctx,
		bucketClient,
		gobFileDir,
		chnSendSDEUpdaterData,
		chnSendDone,
	)

	// Download and write Core .gob files
	// Also, channel-send updater data for use in constants file generation
	chnSendCoreUpdaterData, chnRecvCoreUpdaterData := chanresult.
		NewChanResult[b.CoreUpdaterData](ctx, 1, 0).Split()
	go transceiveDownloadAndWriteCoreBucketData(
		ctx,
		bucketClient,
		gobFileDir,
		chnSendCoreUpdaterData,
		chnSendDone,
	)

	// Download bucket constants data, and receive the 2 updater datas
	constantsData, err := downloadConstantsBucketData(ctx, bucketClient)
	if err != nil {
		return err
	}
	sdeUpdaterData, err := chnRecvSDEUpdaterData.Recv()
	if err != nil {
		return err
	}
	coreUpdaterData, err := chnRecvCoreUpdaterData.Recv()
	if err != nil {
		return err
	}

	// Write constants file using constants bucket data + 2 updater datas
	go transceiveWriteConstants(
		constantsFilePath,
		constantsData,
		sdeUpdaterData,
		coreUpdaterData,
		chnSendDone,
	)

	// Wait for the 3 goroutines to finish
	for i := 0; i < 3; i++ {
		_, err := chnRecvDone.Recv()
		if err != nil {
			return err
		}
	}

	return nil
}
