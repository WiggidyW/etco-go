package bucket

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	SAC "github.com/WiggidyW/etco-go/client/caching/strong/anticaching/single"
)

type WebShopLocationTypeMapsBuilderWriterParams struct {
	WebShopLocationTypeMapsBuilder map[b.TypeId]b.WebShopLocationTypeBundle
}

func (p WebShopLocationTypeMapsBuilderWriterParams) AntiCacheKey() string {
	return cachekeys.WebShopLocationTypeMapsBuilderReaderCacheKey()
}

type SAC_WebShopLocationTypeMapsBuilderWriterClient = SAC.StrongAntiCachingClient[
	WebShopLocationTypeMapsBuilderWriterParams,
	struct{},
	WebShopLocationTypeMapsBuilderWriterClient,
]

func NewSAC_WebShopLocationTypeMapsBuilderWriterClient(
	bucketClient bucket.BucketClient,
	antiCache *cache.StrongAntiCache,
) SAC_WebShopLocationTypeMapsBuilderWriterClient {
	return SAC.NewStrongAntiCachingClient(
		NewWebShopLocationTypeMapsBuilderWriterClient(bucketClient),
		antiCache,
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
