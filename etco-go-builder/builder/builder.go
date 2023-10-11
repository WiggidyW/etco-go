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
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	chnSend, chnRecv := chanresult.
		NewChanResult[struct{}](ctx, 3, 0).Split()

	go transceiveDownloadAndWriteUpdaterBucketData(
		ctx,
		bucketClient,
		constantsFilePath,
		chnSend,
	)
	go transceiveDownloadAndWriteSDEBucketData(
		ctx,
		bucketClient,
		gobFileDir,
		chnSend,
	)
	go transceiveDownloadAndWriteCoreBucketData(
		ctx,
		bucketClient,
		gobFileDir,
		chnSend,
	)

	for i := 0; i < 3; i++ {
		_, err := chnRecv.Recv()
		if err != nil {
			return err
		}
	}

	return nil
}
