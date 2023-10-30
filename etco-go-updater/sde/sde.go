package sde

import (
	"context"
	"net/http"

	b "github.com/WiggidyW/etco-go-bucket"
)

func DownloadAndConvert(
	ctx context.Context,
	httpClient *http.Client,
	userAgent string,
	sdeChecksum string,
	pathSDE string,
) (
	sdeBucketData b.SDEBucketData,
	err error,
) {
	if err := downloadSDE(
		ctx,
		httpClient,
		userAgent,
		pathSDE,
	); err != nil {
		return sdeBucketData, err
	}
	return LoadAndConvert(ctx, sdeChecksum, pathSDE)
}
