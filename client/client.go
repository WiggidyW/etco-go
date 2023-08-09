package client

import (
	"context"
)

type Client[F any, D any] interface {
	Fetch(ctx context.Context, params F) (*D, error)
}

type ClientWrapper[
	FF any,
	F any,
	D any,
	C Client[F, D],
] struct {
	Client C
}
