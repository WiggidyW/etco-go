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
	WEB_S_TYPEMAPSBUILDER_EXPIRES           time.Duration = 24 * time.Hour
	WEB_S_TYPEMAPSBUILDER_MIN_EXPIRES       time.Duration = 0
	WEB_S_TYPEMAPSBUILDER_SLOCK_TTL         time.Duration = 1 * time.Minute
	WEB_S_TYPEMAPSBUILDER_SLOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
)

type WebShopLocationTypeMapsBuilderReaderParams struct{}

func (p WebShopLocationTypeMapsBuilderReaderParams) CacheKey() string {
	return cachekeys.WebShopLocationTypeMapsBuilderReaderCacheKey()
}

type SC_WebShopLocationTypeMapsBuilderReaderClient = SC.StrongCachingClient[
	WebShopLocationTypeMapsBuilderReaderParams,
	map[b.TypeId]b.WebShopLocationTypeBundle,
	cache.ExpirableData[map[b.TypeId]b.WebShopLocationTypeBundle],
	WebShopLocationTypeMapsBuilderReaderClient,
]

func NewSC_WebShopLocationTypeMapsBuilderReaderClient(
	bucketClient bucket.BucketClient,
	sCache cache.SharedServerCache,
) SC_WebShopLocationTypeMapsBuilderReaderClient {
	return SC.NewStrongCachingClient(
		NewWebShopLocationTypeMapsBuilderReaderClient(bucketClient),
		WEB_S_TYPEMAPSBUILDER_MIN_EXPIRES,
		sCache,
		WEB_S_TYPEMAPSBUILDER_SLOCK_TTL,
		WEB_S_TYPEMAPSBUILDER_SLOCK_MAX_BACKOFF,
	)
}

type WebShopLocationTypeMapsBuilderReaderClient struct {
	bucketClient bucket.BucketClient
	expires      time.Duration
}

func NewWebShopLocationTypeMapsBuilderReaderClient(
	bucketClient bucket.BucketClient,
) WebShopLocationTypeMapsBuilderReaderClient {
	return WebShopLocationTypeMapsBuilderReaderClient{
		bucketClient: bucketClient,
		expires:      WEB_S_TYPEMAPSBUILDER_EXPIRES,
	}
}

func (ahsrc WebShopLocationTypeMapsBuilderReaderClient) Fetch(
	ctx context.Context,
	params WebShopLocationTypeMapsBuilderReaderParams,
) (
	rep *cache.ExpirableData[map[b.TypeId]b.WebShopLocationTypeBundle],
	err error,
) {
	d, err := ahsrc.bucketClient.ReadWebShopLocationTypeMapsBuilder(ctx)
	if err != nil {
		return nil, err
	} else {
		return cache.NewExpirableDataPtr(
			d,
			time.Now().Add(ahsrc.expires),
		), nil
	}
}
