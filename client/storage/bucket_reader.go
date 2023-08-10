package storage

import (
	"context"
	"time"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
)

type CachingBucketReaderClient[D any] struct {
	*client.CachingClient[
		BucketReaderParams,
		D,
		cache.ExpirableData[D],
		BucketReaderClient[D],
	]
}

type BucketReaderParams string // object name

func (brp BucketReaderParams) CacheKey() string {
	return string(brp)
}

type BucketReaderClient[D any] struct {
	*BucketClient
	expires time.Duration
}

func (brc BucketReaderClient[D]) Fetch(
	ctx context.Context,
	params BucketReaderParams,
) (*cache.ExpirableData[D], error) {
	d := new(D) // &d
	if _, err := brc.read(ctx, params.CacheKey(), d); err != nil {
		return nil, err
	} else {
		data := cache.NewExpirableData[D](
			*d,
			time.Now().Add(brc.expires),
		)
		return &data, nil
	}
}
