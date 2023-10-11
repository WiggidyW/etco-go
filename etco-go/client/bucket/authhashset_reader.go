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
	AUTH_HASH_SET_EXPIRES        time.Duration = 24 * time.Hour
	AUTH_HASH_SET_MIN_EXPIRES    time.Duration = 0
	AUTH_HASH_SET_SLOCK_TTL      time.Duration = 1 * time.Minute
	AUTH_HASH_SET_SLOCK_MAX_WAIT time.Duration = 1 * time.Minute
)

type AuthHashSetReaderParams struct {
	AuthDomain string
}

func (p AuthHashSetReaderParams) CacheKey() string {
	return cachekeys.AuthHashSetReaderCacheKey(p.AuthDomain)
}

type SC_AuthHashSetReaderClient = SC.StrongCachingClient[
	AuthHashSetReaderParams,
	b.AuthHashSet,
	cache.ExpirableData[b.AuthHashSet],
	AuthHashSetReaderClient,
]

func NewSC_AuthHashSetReaderClient(
	bucketClient bucket.BucketClient,
	sCache cache.SharedServerCache,
) SC_AuthHashSetReaderClient {
	return SC.NewStrongCachingClient(
		NewAuthHashSetReaderClient(bucketClient),
		AUTH_HASH_SET_MIN_EXPIRES,
		sCache,
		AUTH_HASH_SET_SLOCK_TTL,
		AUTH_HASH_SET_SLOCK_MAX_WAIT,
	)
}

type AuthHashSetReaderClient struct {
	bucketClient bucket.BucketClient
	expires      time.Duration
}

func NewAuthHashSetReaderClient(
	bucketClient bucket.BucketClient,
) AuthHashSetReaderClient {
	return AuthHashSetReaderClient{
		bucketClient: bucketClient,
		expires:      AUTH_HASH_SET_EXPIRES,
	}
}

func (ahsrc AuthHashSetReaderClient) Fetch(
	ctx context.Context,
	params AuthHashSetReaderParams,
) (
	rep *cache.ExpirableData[b.AuthHashSet],
	err error,
) {
	ahs, err := ahsrc.bucketClient.ReadAuthHashSet(ctx, params.AuthDomain)
	if err != nil {
		return nil, err
	} else {
		return cache.NewExpirableDataPtr(
			ahs,
			time.Now().Add(ahsrc.expires),
		), nil
	}
}
