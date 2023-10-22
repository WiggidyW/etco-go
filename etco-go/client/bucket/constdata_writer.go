package bucket

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	SAC "github.com/WiggidyW/etco-go/client/caching/strong/anticaching/single"
)

type ConstDataWriterParams struct {
	ConstData b.ConstantsData
}

func (p ConstDataWriterParams) AntiCacheKey() string {
	return cachekeys.ConstDataReaderCacheKey()
}

type SAC_ConstDataWriterClient = SAC.StrongAntiCachingClient[
	ConstDataWriterParams,
	struct{},
	ConstDataWriterClient,
]

func NewSAC_ConstDataWriterClient(
	bucketClient bucket.BucketClient,
	antiCache *cache.StrongAntiCache,
) SAC_ConstDataWriterClient {
	return SAC.NewStrongAntiCachingClient(
		NewConstDataWriterClient(bucketClient),
		antiCache,
	)
}

type ConstDataWriterClient struct {
	bucketClient bucket.BucketClient
}

func NewConstDataWriterClient(
	bucketClient bucket.BucketClient,
) ConstDataWriterClient {
	return ConstDataWriterClient{bucketClient}
}

func (cdwc ConstDataWriterClient) Fetch(
	ctx context.Context,
	params ConstDataWriterParams,
) (
	rep *struct{},
	err error,
) {
	err = cdwc.bucketClient.WriteConstantsData(
		ctx,
		params.ConstData,
	)
	if err != nil {
		return nil, err
	} else {
		return &struct{}{}, nil
	}
}
