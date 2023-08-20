package jwks

import (
	"context"

	"github.com/WiggidyW/weve-esi/cache"
	wc "github.com/WiggidyW/weve-esi/client/caching/weak"
	"github.com/WiggidyW/weve-esi/client/esi/internal/raw"
)

type WC_JWKSClient = wc.WeakCachingClient[
	JWKSParams,
	[]byte,
	cache.ExpirableData[[]byte],
	JWKSClient,
]

type JWKSClient struct {
	raw.RawClient
}

func (jwks JWKSClient) Fetch(
	ctx context.Context,
	params JWKSParams,
) (*cache.ExpirableData[[]byte], error) {
	return jwks.FetchJWKS(ctx)
}
