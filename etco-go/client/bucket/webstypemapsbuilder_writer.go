package bucket

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	SMAC "github.com/WiggidyW/etco-go/client/caching/strong/anticaching/multi"
)

type WebShopLocationTypeMapsBuilderWriterParams struct {
	WebShopLocationTypeMapsBuilder map[b.TypeId]b.WebShopLocationTypeBundle
}

func (p WebShopLocationTypeMapsBuilderWriterParams) AntiCacheKeys() []string {
	return []string{
		cachekeys.WebShopLocationTypeMapsBuilderReaderCacheKey(),
		cachekeys.WebShopBundleKeysCacheKey(),
	}
}

type SAC_WebShopLocationTypeMapsBuilderWriterClient = SMAC.StrongMultiAntiCachingClient[
	WebShopLocationTypeMapsBuilderWriterParams,
	struct{},
	WebShopLocationTypeMapsBuilderWriterClient,
]

func NewSMAC_WebShopLocationTypeMapsBuilderWriterClient(
	bucketClient bucket.BucketClient,
	builderAntiCache *cache.StrongAntiCache,
	bundleKeysAntiCache *cache.StrongAntiCache,
) SAC_WebShopLocationTypeMapsBuilderWriterClient {
	return SMAC.NewStrongMultiAntiCachingClient(
		NewWebShopLocationTypeMapsBuilderWriterClient(bucketClient),
		builderAntiCache,
		bundleKeysAntiCache,
	)
}

type WebShopLocationTypeMapsBuilderWriterClient struct {
	bucketClient bucket.BucketClient
}

func NewWebShopLocationTypeMapsBuilderWriterClient(
	bucketClient bucket.BucketClient,
) WebShopLocationTypeMapsBuilderWriterClient {
	return WebShopLocationTypeMapsBuilderWriterClient{bucketClient}
}

func (wbstmbwc WebShopLocationTypeMapsBuilderWriterClient) Fetch(
	ctx context.Context,
	params WebShopLocationTypeMapsBuilderWriterParams,
) (
	rep *struct{},
	err error,
) {
	err = wbstmbwc.bucketClient.WriteWebShopLocationTypeMapsBuilder(
		ctx,
		params.WebShopLocationTypeMapsBuilder,
	)
	if err != nil {
		return nil, err
	} else {
		return &struct{}{}, nil
	}
}
