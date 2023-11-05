package bucket

import (
	"context"
	"time"
)

const (
	WEB_CAPACITY_MULTIPLIER = 3
	WEB_CAPACITY_DIVISOR    = 2
)

func webGet[K comparable, V any](
	ctx context.Context,
	method func(context.Context, int) (map[K]V, error),
	typeStr, cacheKey string,
	lockTTL, lockMaxBackoff, expiresIn time.Duration,
	makeCap int,
) (
	rep map[K]V,
	expires *time.Time,
	err error,
) {
	return get(
		ctx,
		func(ctx context.Context) (map[K]V, error) {
			return method(ctx, transformWebCapacity(makeCap))
		},
		typeStr, cacheKey,
		lockTTL, lockMaxBackoff, expiresIn,
		makeMapPtrFunc[K, V](makeCap),
	)
}

func makeMapPtrFunc[K comparable, V any](
	capacity int,
) *func() *map[K]V {
	fn := func() *map[K]V {
		m := make(map[K]V, capacity)
		return &m
	}
	return &fn
}

func transformWebCapacity(capacity int) int {
	return capacity * WEB_CAPACITY_MULTIPLIER / WEB_CAPACITY_DIVISOR
}
