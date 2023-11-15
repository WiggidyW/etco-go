package client

import (
	"context"
)

type Client[P any, D any] interface {
	Fetch(ctx context.Context, params P) (*D, error)
}
