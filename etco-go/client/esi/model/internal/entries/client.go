package entries

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/esi/model/internal/model"
	"github.com/WiggidyW/etco-go/client/esi/model/internal/naive"
	"github.com/WiggidyW/etco-go/client/esi/raw_"
)

type EntriesClient[P naive.UrlParams, E any] struct {
	Inner      model.ModelClient[P, []E]
	NumEntries int
}

func NewEntriesClient[P naive.UrlParams, E any](
	rawClient raw_.RawClient,
	numEntries int,
) EntriesClient[P, E] {
	return EntriesClient[P, E]{
		Inner:      model.NewModelClient[P, []E](rawClient),
		NumEntries: numEntries,
	}
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
