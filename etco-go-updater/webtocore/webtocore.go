package webtocore

import (
	"context"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"
)

type WebAndCoreBucketData struct {
	WebBucketData  b.WebBucketData
	CoreBucketData b.CoreBucketData
}

func TransceiveDownloadAndConvert(
	ctx context.Context,
	bucketClient *b.BucketClient,
	chnSend chanresult.ChanSendResult[WebAndCoreBucketData],
) error {
	webAndCoreBucketData, err := DownloadAndConvert(ctx, bucketClient)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(webAndCoreBucketData)
	}
}

func DownloadAndConvert(
	ctx context.Context,
	bucketClient *b.BucketClient,
) (
	webAndCoreBucketData WebAndCoreBucketData,
	err error,
) {
	webBucketData, err := downloadWebBucketData(ctx, bucketClient)
	if err != nil {
		return webAndCoreBucketData, err
	}
	coreBucketData, err := convert(webBucketData)
	if err != nil {
		return webAndCoreBucketData, err
	}
	return WebAndCoreBucketData{
		WebBucketData:  webBucketData,
		CoreBucketData: coreBucketData,
	}, nil
}
