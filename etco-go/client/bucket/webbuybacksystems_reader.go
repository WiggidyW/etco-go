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
	WEB_BUYBACK_SYSTEMS_EXPIRES           time.Duration = 24 * time.Hour
	WEB_BUYBACK_SYSTEMS_MIN_EXPIRES       time.Duration = 0
	WEB_BUYBACK_SYSTEMS_SLOCK_TTL         time.Duration = 1 * time.Minute
	WEB_BUYBACK_SYSTEMS_SLOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
)

type WebBuybackSystemsReaderParams struct{}

func (p WebBuybackSystemsReaderParams) CacheKey() string {
	return cachekeys.WebBuybackSystemsReaderCacheKey()
}

type SC_WebBuybackSystemsReaderClient = SC.StrongCachingClient[
	WebBuybackSystemsReaderParams,
	map[b.SystemId]b.WebBuybackSystem,
	cache.ExpirableData[map[b.SystemId]b.WebBuybackSystem],
	WebBuybackSystemsReaderClient,
]

func NewSC_WebBuybackSystemsReaderClient(
	bucketClient bucket.BucketClient,
	sCache cache.SharedServerCache,
) SC_WebBuybackSystemsReaderClient {
	return SC.NewStrongCachingClient(
		NewWebBuybackSystemsReaderClient(bucketClient),
		WEB_BUYBACK_SYSTEMS_MIN_EXPIRES,
		sCache,
		WEB_BUYBACK_SYSTEMS_SLOCK_TTL,
		WEB_BUYBACK_SYSTEMS_SLOCK_MAX_BACKOFF,
	)
}

type WebBuybackSystemsReaderClient struct {
	bucketClient bucket.BucketClient
	expires      time.Duration
}

func NewWebBuybackSystemsReaderClient(
	bucketClient bucket.BucketClient,
) WebBuybackSystemsReaderClient {
	return WebBuybackSystemsReaderClient{
		bucketClient: bucketClient,
		expires:      WEB_BUYBACK_SYSTEMS_EXPIRES,
	}
}

func (ahsrc WebBuybackSystemsReaderClient) Fetch(
	ctx context.Context,
	params WebBuybackSystemsReaderParams,
) (
	rep *cache.ExpirableData[map[b.SystemId]b.WebBuybackSystem],
	err error,
) {
	d, err := ahsrc.bucketClient.ReadWebBuybackSystems(ctx)
	if err != nil {
		return nil, err
	} else {
		return cache.NewExpirableDataPtr(
			d,
			time.Now().Add(ahsrc.expires),
		), nil
	}
}
