package reader

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	sc "github.com/WiggidyW/etco-go/client/caching/strong/caching"
	b "github.com/WiggidyW/etco-go/client/configure/internal/bucket/internal"
)

type SC_BucketReaderClient[D any] struct {
	sc.StrongCachingClient[
		BucketReaderParams,
		D,
		cache.ExpirableData[D],
		BucketReaderClient[D],
	]
}

type BucketReaderClient[D any] struct {
	*b.BucketClient
	expires time.Duration
}

func (brc BucketReaderClient[D]) Fetch(
	ctx context.Context,
	params BucketReaderParams,
) (*cache.ExpirableData[D], error) {
	d := new(D) // &d
	if _, err := b.Read(
		brc.BucketClient,
		ctx,
		params.CacheKey(),
		d,
	); err != nil {
		return nil, err
	} else {
		data := cache.NewExpirableData[D](
			*d, // possibly == D{} (empty but not nil)
			time.Now().Add(brc.expires),
		)
		return &data, nil
	}
}
