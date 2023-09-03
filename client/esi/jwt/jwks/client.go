package jwks

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	wc "github.com/WiggidyW/etco-go/client/caching/weak"
	"github.com/WiggidyW/etco-go/client/esi/raw_"
)

const (
	JWKS_MIN_EXPIRES    time.Duration = 24 * time.Hour
	JWKS_SLOCK_TTL      time.Duration = 30 * time.Second
	JWKS_SLOCK_MAX_WAIT time.Duration = 10 * time.Second
)

type WC_JWKSClient = wc.WeakCachingClient[
	JWKSParams,
	[]byte,
	cache.ExpirableData[[]byte],
	JWKSClient,
]

func NewWC_JWKSClient(
	rawClient raw_.RawClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) WC_JWKSClient {
	return wc.NewWeakCachingClient[
		JWKSParams,
		[]byte,
		cache.ExpirableData[[]byte],
		JWKSClient,
	](
		NewJWKSClient(rawClient),
		JWKS_MIN_EXPIRES,
		cCache,
		sCache,
		JWKS_SLOCK_TTL,
		JWKS_SLOCK_MAX_WAIT,
	)
}

type JWKSClient struct {
	raw_.RawClient
}

func NewJWKSClient(rawClient raw_.RawClient) JWKSClient {
	return JWKSClient{
		RawClient: rawClient,
	}
}

func (jwks JWKSClient) Fetch(
	ctx context.Context,
	params JWKSParams,
) (*cache.ExpirableData[[]byte], error) {
	return jwks.FetchJWKS(ctx)
}
