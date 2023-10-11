package sde

import (
	"context"
	"net/http"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"
)

func TransceiveDownloadAndConvert(
	ctx context.Context,
	httpClient *http.Client,
	userAgent string,
	pathSDE string,
	chnSend chanresult.ChanSendResult[b.SDEBucketData],
) error {
	etcoSDEBucketData, err := DownloadAndConvert(
		ctx,
		httpClient,
		userAgent,
		pathSDE,
	)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(etcoSDEBucketData)
	}
}

func DownloadAndConvert(
	ctx context.Context,
	httpClient *http.Client,
	userAgent string,
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
	return LoadAndConvert(ctx, pathSDE)
}
