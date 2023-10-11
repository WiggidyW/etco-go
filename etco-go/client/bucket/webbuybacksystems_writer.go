package bucket

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	SAC "github.com/WiggidyW/etco-go/client/caching/strong/anticaching/single"
)

type WebBuybackSystemsWriterParams struct {
	WebBuybackSystems map[b.SystemId]b.WebBuybackSystem
}

func (p WebBuybackSystemsWriterParams) AntiCacheKey() string {
	return cachekeys.WebBuybackSystemsReaderCacheKey()
}

type SAC_WebBuybackSystemsWriterClient = SAC.StrongAntiCachingClient[
	WebBuybackSystemsWriterParams,
	struct{},
	WebBuybackSystemsWriterClient,
]

func NewSAC_WebBuybackSystemsWriterClient(
	bucketClient bucket.BucketClient,
	antiCache *cache.StrongAntiCache,
) SAC_WebBuybackSystemsWriterClient {
	return SAC.NewStrongAntiCachingClient(
		NewWebBuybackSystemsWriterClient(bucketClient),
		antiCache,
	)
}

type WebBuybackSystemsWriterClient struct {
	bucketClient bucket.BucketClient
}

func NewWebBuybackSystemsWriterClient(
	bucketClient bucket.BucketClient,
) WebBuybackSystemsWriterClient {
	return WebBuybackSystemsWriterClient{
		bucketClient: bucketClient,
	}
}

func (wbstmbwc WebBuybackSystemsWriterClient) Fetch(
	ctx context.Context,
	params WebBuybackSystemsWriterParams,
) (
	rep *struct{},
	err error,
) {
	err = wbstmbwc.bucketClient.WriteWebBuybackSystems(
		ctx,
		params.WebBuybackSystems,
	)
	if err != nil {
		return nil, err
	} else {
		return &struct{}{}, nil
	}
}
