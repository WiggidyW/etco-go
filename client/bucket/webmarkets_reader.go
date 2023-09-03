package bucket

import (
	"context"
	"time"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	SC "github.com/WiggidyW/etco-go/client/caching/strong/caching"
)

const (
	WEB_MARKETS_EXPIRES        time.Duration = 24 * time.Hour
	WEB_MARKETS_MIN_EXPIRES    time.Duration = 0
	WEB_MARKETS_SLOCK_TTL      time.Duration = 1 * time.Minute
	WEB_MARKETS_SLOCK_MAX_WAIT time.Duration = 1 * time.Minute
)

type WebMarketsReaderParams struct{}

func (p WebMarketsReaderParams) CacheKey() string {
	return cachekeys.WebMarketsReaderCacheKey()
}

type SC_WebMarketsReaderClient = SC.StrongCachingClient[
	WebMarketsReaderParams,
	map[b.MarketName]b.WebMarket,
	cache.ExpirableData[map[b.MarketName]b.WebMarket],
	WebMarketsReaderClient,
]

func NewSC_WebMarketsReaderClient(
	bucketClient bucket.BucketClient,
	sCache cache.SharedServerCache,
) SC_WebMarketsReaderClient {
	return SC.NewStrongCachingClient(
		NewWebMarketsReaderClient(bucketClient),
		WEB_MARKETS_MIN_EXPIRES,
		sCache,
		WEB_MARKETS_SLOCK_TTL,
		WEB_MARKETS_SLOCK_MAX_WAIT,
	)
}

type WebMarketsReaderClient struct {
	bucketClient bucket.BucketClient
	expires      time.Duration
}

func NewWebMarketsReaderClient(
	bucketClient bucket.BucketClient,
) WebMarketsReaderClient {
	return WebMarketsReaderClient{
		bucketClient: bucketClient,
		expires:      WEB_MARKETS_EXPIRES,
	}
}

func (ahsrc WebMarketsReaderClient) Fetch(
	ctx context.Context,
	params WebMarketsReaderParams,
) (
	rep *cache.ExpirableData[map[b.MarketName]b.WebMarket],
	err error,
) {
	d, err := ahsrc.bucketClient.ReadWebMarkets(ctx)
	if err != nil {
		return nil, err
	} else {
		return cache.NewExpirableDataPtr(
			d,
			time.Now().Add(ahsrc.expires),
		), nil
	}
}
