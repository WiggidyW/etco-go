package webtocore

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"
)

type WebAndCoreBucketData struct {
	WebBucketData  b.WebBucketData
	CoreBucketData b.CoreBucketData
}

func DownloadAndConvert(
	ctx context.Context,
	bucketClient *b.BucketClient,
	webAttrs WebAttrs,
	sdeSystems map[b.SystemId]b.System,
) (
	coreBucketData b.CoreBucketData,
	err error,
) {
	webBucketData, err := downloadWebBucketData(ctx, bucketClient)
	if err != nil {
		return coreBucketData, err
	}
	return convert(webBucketData, webAttrs, sdeSystems)
}
