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
	CONSTDATA_EXPIRES           time.Duration = 24 * time.Hour
	CONSTDATA_MIN_EXPIRES       time.Duration = 0
	CONSTDATA_SLOCK_TTL         time.Duration = 1 * time.Minute
	CONSTDATA_SLOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
)

type ConstDataReaderParams struct{}

func (p ConstDataReaderParams) CacheKey() string {
	return cachekeys.ConstDataReaderCacheKey()
}

type SC_ConstDataReaderClient = SC.StrongCachingClient[
	ConstDataReaderParams,
	b.ConstantsData,
	cache.ExpirableData[b.ConstantsData],
	ConstDataReaderClient,
]

func NewSC_ConstDataReaderClient(
	bucketClient bucket.BucketClient,
	sCache cache.SharedServerCache,
) SC_ConstDataReaderClient {
	return SC.NewStrongCachingClient(
		NewConstDataReaderClient(bucketClient),
		CONSTDATA_MIN_EXPIRES,
		sCache,
		CONSTDATA_SLOCK_TTL,
		CONSTDATA_SLOCK_MAX_BACKOFF,
	)
}

type ConstDataReaderClient struct {
	bucketClient bucket.BucketClient
	expires      time.Duration
}

func NewConstDataReaderClient(
	bucketClient bucket.BucketClient,
) ConstDataReaderClient {
	return ConstDataReaderClient{
		bucketClient: bucketClient,
		expires:      CONSTDATA_EXPIRES,
	}
}

func (cdrc ConstDataReaderClient) Fetch(
	ctx context.Context,
	params ConstDataReaderParams,
) (
	rep *cache.ExpirableData[b.ConstantsData],
	err error,
) {
	d, err := cdrc.bucketClient.ReadConstantsData(ctx)
	if err != nil {
		return nil, err
	} else {
		return cache.NewExpirableDataPtr(
			d,
			time.Now().Add(cdrc.expires),
		), nil
	}
}
