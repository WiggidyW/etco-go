package staticdb

type Container[T any] struct {
	Inner T
}

func NewContainer[T any](inner T) *Container[T] {
	return &Container[T]{Inner: inner}
}
