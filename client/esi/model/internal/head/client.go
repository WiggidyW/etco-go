package head

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/cache"
	"github.com/WiggidyW/eve-trading-co-go/client/esi/model/internal/naive"
)

type HeadClient[P naive.UrlParams] struct {
	naive.NaiveClient[P]
}

func (hc HeadClient[P]) Fetch(
	ctx context.Context,
	params naive.NaiveParams[P],
) (*cache.ExpirableData[int], error) {
	return hc.FetchHead(ctx, params)
}
