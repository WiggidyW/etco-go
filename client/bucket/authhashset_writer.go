package bucket

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	SAC "github.com/WiggidyW/etco-go/client/caching/strong/anticaching/single"
)

type AuthHashSetWriterParams struct {
	AuthDomain  string
	AuthHashSet b.AuthHashSet
}

func (p AuthHashSetWriterParams) AntiCacheKey() string {
	return cachekeys.AuthHashSetReaderCacheKey(p.AuthDomain)
}

type SAC_AuthHashSetWriterClient = SAC.StrongAntiCachingClient[
	AuthHashSetWriterParams,
	struct{},
	AuthHashSetWriterClient,
]

func NewSAC_AuthHashSetWriterClient(
	bucketClient bucket.BucketClient,
	antiCache *cache.StrongAntiCache,
) SAC_AuthHashSetWriterClient {
	return SAC.NewStrongAntiCachingClient(
		AuthHashSetWriterClient{bucketClient},
		antiCache,
	)
}

type AuthHashSetWriterClient struct {
	bucketClient bucket.BucketClient
}

func (ahsrc AuthHashSetWriterClient) Fetch(
	ctx context.Context,
	params AuthHashSetWriterParams,
) (
	rep *struct{},
	err error,
) {
	err = ahsrc.bucketClient.WriteAuthHashSet(
		ctx,
		params.AuthHashSet,
		params.AuthDomain,
	)
	if err != nil {
		return nil, err
	} else {
		return &struct{}{}, nil
	}
}
