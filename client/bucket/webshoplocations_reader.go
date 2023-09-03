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
	WEB_SHOP_LOCATIONS_EXPIRES        time.Duration = 24 * time.Hour
	WEB_SHOP_LOCATIONS_MIN_EXPIRES    time.Duration = 0
	WEB_SHOP_LOCATIONS_SLOCK_TTL      time.Duration = 1 * time.Minute
	WEB_SHOP_LOCATIONS_SLOCK_MAX_WAIT time.Duration = 1 * time.Minute
)

type WebShopLocationsReaderParams struct{}

func (p WebShopLocationsReaderParams) CacheKey() string {
	return cachekeys.WebShopLocationsReaderCacheKey()
}

type SC_WebShopLocationsReaderClient = SC.StrongCachingClient[
	WebShopLocationsReaderParams,
	map[b.LocationId]b.WebShopLocation,
	cache.ExpirableData[map[b.LocationId]b.WebShopLocation],
	WebShopLocationsReaderClient,
]

func NewSC_WebShopLocationsReaderClient(
	bucketClient bucket.BucketClient,
	sCache cache.SharedServerCache,
) SC_WebShopLocationsReaderClient {
	return SC.NewStrongCachingClient(
		NewWebShopLocationsReaderClient(bucketClient),
		WEB_SHOP_LOCATIONS_MIN_EXPIRES,
		sCache,
		WEB_SHOP_LOCATIONS_SLOCK_TTL,
		WEB_SHOP_LOCATIONS_SLOCK_MAX_WAIT,
	)
}

type WebShopLocationsReaderClient struct {
	bucketClient bucket.BucketClient
	expires      time.Duration
}

func NewWebShopLocationsReaderClient(
	bucketClient bucket.BucketClient,
) WebShopLocationsReaderClient {
	return WebShopLocationsReaderClient{
		bucketClient: bucketClient,
		expires:      WEB_SHOP_LOCATIONS_EXPIRES,
	}
}

func (ahsrc WebShopLocationsReaderClient) Fetch(
	ctx context.Context,
	params WebShopLocationsReaderParams,
) (
	rep *cache.ExpirableData[map[b.LocationId]b.WebShopLocation],
	err error,
) {
	d, err := ahsrc.bucketClient.ReadWebShopLocations(ctx)
	if err != nil {
		return nil, err
	} else {
		return cache.NewExpirableDataPtr(
			d,
			time.Now().Add(ahsrc.expires),
		), nil
	}
}
