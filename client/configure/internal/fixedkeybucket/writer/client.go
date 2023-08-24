package writer

import (
	"context"

	bucketwriter "github.com/WiggidyW/eve-trading-co-go/client/configure/internal/bucket/writer"
)

type SAC_FixedKeyBucketWriterClient[D any] struct {
	Inner      bucketwriter.SAC_BucketWriterClient[D]
	ObjectName string
}

func (fkbwc SAC_FixedKeyBucketWriterClient[D]) Fetch(
	ctx context.Context,
	params FixedKeyBucketWriterParams[D],
) (*struct{}, error) {
	return fkbwc.Inner.Fetch(ctx, bucketwriter.BucketWriterParams[D]{
		ObjectName: fkbwc.ObjectName,
		Val:        params.Val,
	})
}

type FixedKeyBucketWriterClient[D any] struct {
	Inner      bucketwriter.BucketWriterClient[D]
	ObjectName string
}

func (fkbwc FixedKeyBucketWriterClient[D]) Fetch(
	ctx context.Context,
	params FixedKeyBucketWriterParams[D],
) (*struct{}, error) {
	return fkbwc.Inner.Fetch(ctx, bucketwriter.BucketWriterParams[D]{
		ObjectName: fkbwc.ObjectName,
		Val:        params.Val,
	})
}
