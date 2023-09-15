package bucket

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	SMAC "github.com/WiggidyW/etco-go/client/caching/strong/anticaching/multi"
)

type WebBuybackSystemTypeMapsBuilderWriterParams struct {
	WebBuybackSystemTypeMapsBuilder map[b.TypeId]b.WebBuybackSystemTypeBundle
}

func (p WebBuybackSystemTypeMapsBuilderWriterParams) AntiCacheKeys() []string {
	return []string{
		cachekeys.WebBuybackSystemTypeMapsBuilderReaderCacheKey(),
		cachekeys.WebBuybackBundleKeysCacheKey(),
	}
}

type SAC_WebBuybackSystemTypeMapsBuilderWriterClient = SMAC.StrongMultiAntiCachingClient[
	WebBuybackSystemTypeMapsBuilderWriterParams,
	struct{},
	WebBuybackSystemTypeMapsBuilderWriterClient,
]

func NewSMAC_WebBuybackSystemTypeMapsBuilderWriterClient(
	bucketClient bucket.BucketClient,
	builderAntiCache *cache.StrongAntiCache,
	bundleKeysAntiCache *cache.StrongAntiCache,
) SAC_WebBuybackSystemTypeMapsBuilderWriterClient {
	return SMAC.NewStrongMultiAntiCachingClient(
		NewWebBuybackSystemTypeMapsBuilderWriterClient(bucketClient),
		builderAntiCache,
		bundleKeysAntiCache,
	)
}

type WebBuybackSystemTypeMapsBuilderWriterClient struct {
	bucketClient bucket.BucketClient
}

func NewWebBuybackSystemTypeMapsBuilderWriterClient(
	bucketClient bucket.BucketClient,
) (
	client WebBuybackSystemTypeMapsBuilderWriterClient,
) {
	return WebBuybackSystemTypeMapsBuilderWriterClient{
		bucketClient: bucketClient,
	}
}

func (wbstmbwc WebBuybackSystemTypeMapsBuilderWriterClient) Fetch(
	ctx context.Context,
	params WebBuybackSystemTypeMapsBuilderWriterParams,
) (
	rep *struct{},
	err error,
) {
	err = wbstmbwc.bucketClient.WriteWebBuybackSystemTypeMapsBuilder(
		ctx,
		params.WebBuybackSystemTypeMapsBuilder,
	)
	if err != nil {
		return nil, err
	} else {
		return &struct{}{}, nil
	}
}
