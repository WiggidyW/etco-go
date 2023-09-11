package bucket

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	SMAC "github.com/WiggidyW/etco-go/client/caching/strong/anticaching/multi"
)

type WebMarketsWriterParams struct {
	WebMarkets map[b.MarketName]b.WebMarket
}

func (p WebMarketsWriterParams) AntiCacheKeys() []string {
	return []string{
		cachekeys.WebMarketsReaderCacheKey(),
		cachekeys.WebMarketsNamesCacheKey(),
	}
}

type SAC_WebMarketsWriterClient = SMAC.StrongMultiAntiCachingClient[
	WebMarketsWriterParams,
	struct{},
	WebMarketsWriterClient,
]

func NewSMAC_WebMarketsWriterClient(
	bucketClient bucket.BucketClient,
	antiCache *cache.StrongAntiCache,
) SAC_WebMarketsWriterClient {
	return SMAC.NewStrongMultiAntiCachingClient(
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
