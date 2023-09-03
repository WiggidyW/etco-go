package bucket

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	SAC "github.com/WiggidyW/etco-go/client/caching/strong/anticaching/single"
)

type WebBuybackSystemTypeMapsBuilderWriterParams struct {
	WebBuybackSystemTypeMapsBuilder map[b.TypeId]b.WebBuybackSystemTypeBundle
}

func (p WebBuybackSystemTypeMapsBuilderWriterParams) AntiCacheKey() string {
	return cachekeys.WebBuybackSystemTypeMapsBuilderReaderCacheKey()
}

type SAC_WebBuybackSystemTypeMapsBuilderWriterClient = SAC.StrongAntiCachingClient[
	WebBuybackSystemTypeMapsBuilderWriterParams,
	struct{},
	WebBuybackSystemTypeMapsBuilderWriterClient,
]

func NewSAC_WebBuybackSystemTypeMapsBuilderWriterClient(
	bucketClient bucket.BucketClient,
	antiCache *cache.StrongAntiCache,
) SAC_WebBuybackSystemTypeMapsBuilderWriterClient {
	return SAC.NewStrongAntiCachingClient(
		NewWebBuybackSystemTypeMapsBuilderWriterClient(bucketClient),
		antiCache,
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
