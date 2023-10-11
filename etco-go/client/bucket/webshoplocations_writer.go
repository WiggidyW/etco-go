package bucket

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	SAC "github.com/WiggidyW/etco-go/client/caching/strong/anticaching/single"
)

type WebShopLocationsWriterParams struct {
	WebShopLocations map[b.LocationId]b.WebShopLocation
}

func (p WebShopLocationsWriterParams) AntiCacheKey() string {
	return cachekeys.WebShopLocationsReaderCacheKey()
}

type SAC_WebShopLocationsWriterClient = SAC.StrongAntiCachingClient[
	WebShopLocationsWriterParams,
	struct{},
	WebShopLocationsWriterClient,
]

func NewSAC_WebShopLocationsWriterClient(
	bucketClient bucket.BucketClient,
	antiCache *cache.StrongAntiCache,
) SAC_WebShopLocationsWriterClient {
	return SAC.NewStrongAntiCachingClient(
		NewWebShopLocationsWriterClient(bucketClient),
		antiCache,
	)
}

type WebShopLocationsWriterClient struct {
	bucketClient bucket.BucketClient
}

func NewWebShopLocationsWriterClient(
	bucketClient bucket.BucketClient,
) WebShopLocationsWriterClient {
	return WebShopLocationsWriterClient{bucketClient}
}

func (wbstmbwc WebShopLocationsWriterClient) Fetch(
	ctx context.Context,
	params WebShopLocationsWriterParams,
) (
	rep *struct{},
	err error,
) {
	err = wbstmbwc.bucketClient.WriteWebShopLocations(
		ctx,
		params.WebShopLocations,
	)
	if err != nil {
		return nil, err
	} else {
		return &struct{}{}, nil
	}
}
