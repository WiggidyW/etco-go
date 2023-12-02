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
	chnSendModified chanresult.ChanSendResult[*b.SDEBucketData],
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
	modified *b.SDEBucketData, // Some if modified, nil if not
	err error,
) {
	if skipSde {
		return nil, nil
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
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	)
	if err != nil {
		return nil, err
	}

	// receive the SDE checksum
	sdeChecksum, err := chnRecvChecksum.Recv()
	if err != nil {
		return nil, err
	} else if sdeChecksum == prevSDEBucketData.UpdaterData.CHECKSUM_SDE {
		// return false if checksums match
		return nil, nil
	}

	// download, convert, and write SDE build bucket data
	var sdeBucketData b.SDEBucketData
	sdeBucketData, err = updateSDE(
		ctx,
		bucketClient,
		httpClient,
		userAgent,
		sdeChecksum,
	)
	if err != nil {
		return nil, err
	}

	return &sdeBucketData, nil
}

func updateSDE(
	ctx context.Context,
	bucketClient *b.BucketClient,
	httpClient *http.Client,
	userAgent string,
	sdeChecksum string,
) (
	sdeBucketData b.SDEBucketData,
	err error,
) {
	// create a temporary directory
	tempDir, err := os.MkdirTemp("", "etco-go-updater-*")
	if err != nil {
		return sdeBucketData, err
	}
	defer os.RemoveAll(tempDir)

	// download and convert the SDE into writeable bucket data
	sdeBucketData, err = updatersde.DownloadAndConvert(
		ctx,
		httpClient,
		userAgent,
		sdeChecksum,
		tempDir,
	)
	if err != nil {
		return sdeBucketData, err
	}

	// write the data to the bucket
	err = bucketClient.WriteSDEData(ctx, sdeBucketData)
	return sdeBucketData, err
}
