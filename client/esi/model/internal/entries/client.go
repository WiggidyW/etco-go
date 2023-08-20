package entries

import (
	"context"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client/esi/model/internal/model"
	"github.com/WiggidyW/weve-esi/client/esi/model/internal/naive"
)

type EntriesClient[P naive.UrlParams, E any] struct {
	Inner      model.ModelClient[P, []E]
	NumEntries int
}

func (mc EntriesClient[P, E]) Fetch(
	ctx context.Context,
	params naive.NaiveParams[P],
) (*cache.ExpirableData[[]E], error) {
	return mc.Inner.Fetch(
		ctx,
		model.ModelParams[P, []E]{
			NaiveParams: params,
			Model:       mc.makeModel(),
		},
	)
}

func (mc EntriesClient[P, E]) makeModel() *[]E {
	model := make([]E, 0, mc.NumEntries)
	return &model
}
