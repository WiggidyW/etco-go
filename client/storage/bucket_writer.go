package storage

import (
	"context"

	"github.com/WiggidyW/weve-esi/client"
)

type AntiCachingBucketWriterClient[D any] struct {
	*client.AntiCachingClient[
		BucketWriterParams[D],
		struct{},
		D,
		BucketWriterClient[D],
	]
}

type BucketWriterParams[D any] struct {
	Key string // object name (domain key + access type)
	Val D
}

func (bwp BucketWriterParams[D]) AntiCacheKey() string {
	return bwp.Key
}

type BucketWriterClient[D any] struct {
	*BucketClient
}

func (bwc BucketWriterClient[D]) Fetch(
	ctx context.Context,
	params BucketWriterParams[D],
) (*struct{}, error) {
	if err := bwc.write(ctx, params.Key, params.Val); err != nil {
		return nil, err
	} else {
		return &struct{}{}, nil
	}
}
