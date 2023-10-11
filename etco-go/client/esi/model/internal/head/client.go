package head

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/esi/model/internal/naive"
	"github.com/WiggidyW/etco-go/client/esi/raw_"
)

type HeadClient[P naive.UrlParams] struct {
	naive.NaiveClient[P]
}

func NewHeadClient[P naive.UrlParams](rawClient raw_.RawClient) HeadClient[P] {
	return HeadClient[P]{
		NaiveClient: naive.NaiveClient[P](rawClient),
	}
}

func (hc HeadClient[P]) Fetch(
	ctx context.Context,
	params naive.NaiveParams[P],
) (*cache.ExpirableData[int], error) {
	return hc.FetchHead(ctx, params)
}
