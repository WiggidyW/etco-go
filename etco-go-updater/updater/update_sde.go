package updater

import (
	"context"
	"net/http"
	"os"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	updatersde "github.com/WiggidyW/etco-go-updater/sde"
)

func transceiveUpdateSDEIfModified(
	ctx context.Context,
	bucketClient *b.BucketClient,
	httpClient *http.Client,
	userAgent string,
	skipSde bool,
	chnSendModified chanresult.ChanSendResult[bool],
) error {
	modified, err := UpdateSDEIfModified(
		ctx,
		bucketClient,
		httpClient,
		userAgent,
		skipSde,
	)
	if err != nil {
		return chnSendModified.SendErr(err)
	} else {
		return chnSendModified.SendOk(modified)
	}
}
func UpdateSDEIfModified(
	ctx context.Context,
	bucketClient *b.BucketClient,
	httpClient *http.Client,
	userAgent string,
	skipSde bool,
) (
	modified bool,
	err error,
) {
	if skipSde {
		return false, nil
	}

	// fetch the SDE checksum in a goroutine
	chnSendChecksum, chnRecvChecksum :=
		chanresult.NewChanResult[string](ctx, 1, 0).Split()
	go updatersde.TransceiveDownloadChecksum(
		ctx,
		httpClient,
		userAgent,
		skipSde,
		chnSendChecksum,
	)

	// fetch previous SDE data
	prevSDEBucketData, err := bucketClient.ReadSDEData(
		ctx,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	)
	if err != nil {
		return false, err
	}

	// receive the SDE checksum
	sdeChecksum, err := chnRecvChecksum.Recv()
	if err != nil {
		return false, err
	} else if sdeChecksum == prevSDEBucketData.UpdaterData.CHECKSUM_SDE {
		// return false if checksums match
		return false, nil
	}

	// download, convert, and write SDE build bucket data
	err = updateSDE(
		ctx,
		bucketClient,
		httpClient,
		userAgent,
		sdeChecksum,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func updateSDE(
	ctx context.Context,
	bucketClient *b.BucketClient,
	httpClient *http.Client,
	userAgent string,
	sdeChecksum string,
) error {
	// create a temporary directory
	tempDir, err := os.MkdirTemp("", "etco-go-updater-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// download and convert the SDE into writeable bucket data
	sdeBucketData, err := updatersde.DownloadAndConvert(
		ctx,
		httpClient,
		userAgent,
		sdeChecksum,
		tempDir,
	)
	if err != nil {
		return err
	}

	// write the data to the bucket
	return bucketClient.WriteSDEData(ctx, sdeBucketData)
}
