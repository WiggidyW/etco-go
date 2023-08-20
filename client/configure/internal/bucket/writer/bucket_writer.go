package writer

import (
	"context"

	sac "github.com/WiggidyW/weve-esi/client/caching/strong/anticaching/single"
	b "github.com/WiggidyW/weve-esi/client/configure/internal/bucket/internal"
)

type SAC_BucketWriterClient[D any] struct {
	*sac.StrongAntiCachingClient[
		BucketWriterParams[D],
		struct{},
		BucketWriterClient[D],
	]
}

type BucketWriterClient[D any] struct {
	*b.BucketClient
}

func (bwc BucketWriterClient[D]) Fetch(
	ctx context.Context,
	params BucketWriterParams[D],
) (*struct{}, error) {
	if err := bwc.Write(ctx, params.ObjectName, params.Val); err != nil {
		return nil, err
	} else {
		return &struct{}{}, nil
	}
}
