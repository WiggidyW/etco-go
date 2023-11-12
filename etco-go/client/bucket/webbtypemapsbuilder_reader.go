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
	WEB_B_TYPEMAPSBUILDER_EXPIRES           time.Duration = 24 * time.Hour
	WEB_B_TYPEMAPSBUILDER_MIN_EXPIRES       time.Duration = 0
	WEB_B_TYPEMAPSBUILDER_SLOCK_TTL         time.Duration = 1 * time.Minute
	WEB_B_TYPEMAPSBUILDER_SLOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
)

type WebBuybackSystemTypeMapsBuilderReaderParams struct{}

func (p WebBuybackSystemTypeMapsBuilderReaderParams) CacheKey() string {
	return cachekeys.WebBuybackSystemTypeMapsBuilderReaderCacheKey()
}

type SC_WebBuybackSystemTypeMapsBuilderReaderClient = SC.StrongCachingClient[
	WebBuybackSystemTypeMapsBuilderReaderParams,
	map[b.TypeId]b.WebBuybackSystemTypeBundle,
	cache.ExpirableData[map[b.TypeId]b.WebBuybackSystemTypeBundle],
	WebBuybackSystemTypeMapsBuilderReaderClient,
]

func NewSC_WebBuybackSystemTypeMapsBuilderReaderClient(
	bucketClient bucket.BucketClient,
	sCache cache.SharedServerCache,
) SC_WebBuybackSystemTypeMapsBuilderReaderClient {
	return SC.NewStrongCachingClient(
		NewWebBuybackSystemTypeMapsBuilderReaderClient(bucketClient),
		WEB_B_TYPEMAPSBUILDER_MIN_EXPIRES,
		sCache,
		WEB_B_TYPEMAPSBUILDER_SLOCK_TTL,
		WEB_B_TYPEMAPSBUILDER_SLOCK_MAX_BACKOFF,
	)
}

type WebBuybackSystemTypeMapsBuilderReaderClient struct {
	bucketClient bucket.BucketClient
	expires      time.Duration
}

func NewWebBuybackSystemTypeMapsBuilderReaderClient(
	bucketClient bucket.BucketClient,
) WebBuybackSystemTypeMapsBuilderReaderClient {
	return WebBuybackSystemTypeMapsBuilderReaderClient{
		bucketClient: bucketClient,
		expires:      WEB_B_TYPEMAPSBUILDER_EXPIRES,
	}
}

func (ahsrc WebBuybackSystemTypeMapsBuilderReaderClient) Fetch(
	ctx context.Context,
	params WebBuybackSystemTypeMapsBuilderReaderParams,
) (
	rep *cache.ExpirableData[map[b.TypeId]b.WebBuybackSystemTypeBundle],
	err error,
) {
	d, err := ahsrc.bucketClient.ReadWebBuybackSystemTypeMapsBuilder(ctx)
	if err != nil {
		return nil, err
	} else {
		return cache.NewExpirableDataPtr(
			d,
			time.Now().Add(ahsrc.expires),
		), nil
	}
}
