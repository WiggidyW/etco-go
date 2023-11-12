package bucket

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
)

const (
	WEB_CAPACITY_MULTIPLIER = 3
	WEB_CAPACITY_DIVISOR    = 2
)

func webGet[K comparable, V any](
	x cache.Context,
	method func(context.Context, int) (map[K]V, error),
	cacheKey, typeStr string,
	expiresIn time.Duration,
	makeCap int,
) (
	rep map[K]V,
	expires time.Time,
	err error,
) {
	return get(
		x,
		func(ctx context.Context) (map[K]V, error) {
			return method(ctx, transformWebCapacity(makeCap))
		},
		cacheKey, typeStr,
		expiresIn,
		makeMapPtrFunc[K, V](makeCap),
	)
}

func makeMapPtrFunc[K comparable, V any](
	capacity int,
) func() *map[K]V {
	return func() *map[K]V {
		m := make(map[K]V, capacity)
		return &m
	}
}

func transformWebCapacity(capacity int) int {
	return capacity * WEB_CAPACITY_MULTIPLIER / WEB_CAPACITY_DIVISOR
}
