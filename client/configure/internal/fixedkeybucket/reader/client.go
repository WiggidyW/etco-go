package reader

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/cache"
	"github.com/WiggidyW/eve-trading-co-go/client/caching"
	bucketreader "github.com/WiggidyW/eve-trading-co-go/client/configure/internal/bucket/reader"
)

type SC_FixedKeyBucketReaderClient[D any] struct {
	Inner      bucketreader.SC_BucketReaderClient[D]
	ObjectName string
}

func (scfkbrc SC_FixedKeyBucketReaderClient[D]) Fetch(
	ctx context.Context,
	params struct{},
) (*caching.CachingResponse[D], error) {
	return scfkbrc.Inner.Fetch(ctx, bucketreader.BucketReaderParams{
		ObjectName: scfkbrc.ObjectName,
	})
}

type FixedKeyBucketReaderClient[D any] struct {
	Inner      bucketreader.BucketReaderClient[D]
	ObjectName string
}

func (fkbrc FixedKeyBucketReaderClient[D]) Fetch(
	ctx context.Context,
	params struct{},
) (*cache.ExpirableData[D], error) {
	return fkbrc.Inner.Fetch(ctx, bucketreader.BucketReaderParams{
		ObjectName: fkbrc.ObjectName,
	})
}
