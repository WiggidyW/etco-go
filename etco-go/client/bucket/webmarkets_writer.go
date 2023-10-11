package bucket

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	SAC "github.com/WiggidyW/etco-go/client/caching/strong/anticaching/single"
)

type WebMarketsWriterParams struct {
	WebMarkets map[b.MarketName]b.WebMarket
}

func (p WebMarketsWriterParams) AntiCacheKey() string {
	return cachekeys.WebMarketsReaderCacheKey()
}

type SAC_WebMarketsWriterClient = SAC.StrongAntiCachingClient[
	WebMarketsWriterParams,
	struct{},
	WebMarketsWriterClient,
]

func NewSAC_WebMarketsWriterClient(
	bucketClient bucket.BucketClient,
	antiCache *cache.StrongAntiCache,
) SAC_WebMarketsWriterClient {
	return SAC.NewStrongAntiCachingClient(
		NewWebMarketsWriterClient(bucketClient),
		antiCache,
	)
}

type WebMarketsWriterClient struct {
	bucketClient bucket.BucketClient
}

func NewWebMarketsWriterClient(
	bucketClient bucket.BucketClient,
) WebMarketsWriterClient {
	return WebMarketsWriterClient{bucketClient}
}

func (wbstmbwc WebMarketsWriterClient) Fetch(
	ctx context.Context,
	params WebMarketsWriterParams,
) (
	rep *struct{},
	err error,
) {
	err = wbstmbwc.bucketClient.WriteWebMarkets(
		ctx,
		params.WebMarkets,
	)
	if err != nil {
		return nil, err
	} else {
		return &struct{}{}, nil
	}
}
