package model

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/cache"
	"github.com/WiggidyW/eve-trading-co-go/client/esi/model/internal/naive"
)

type ModelClient[P naive.UrlParams, M any] struct {
	naive.NaiveClient[P]
}

func (mc ModelClient[P, M]) Fetch(
	ctx context.Context,
	params ModelParams[P, M],
) (*cache.ExpirableData[M], error) {
	return naive.FetchModel[P, M](
		mc.NaiveClient,
		ctx,
		params.NaiveParams,
		params.Model,
	)
}
