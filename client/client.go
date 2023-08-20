package client

import (
	"context"
)

type Client[F any, D any] interface {
	Fetch(ctx context.Context, params F) (*D, error)
}
