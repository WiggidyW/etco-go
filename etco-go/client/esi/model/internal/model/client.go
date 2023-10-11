package model

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/esi/model/internal/naive"
	"github.com/WiggidyW/etco-go/client/esi/raw_"
)

type ModelClient[P naive.UrlParams, M any] struct {
	naive.NaiveClient[P]
}

func NewModelClient[P naive.UrlParams, M any](
	rawClient raw_.RawClient,
) ModelClient[P, M] {
	return ModelClient[P, M]{
		NaiveClient: naive.NaiveClient[P](rawClient),
	}
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
