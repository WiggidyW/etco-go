package jwks

import "github.com/WiggidyW/eve-trading-co-go/client/cachekeys"

type JWKSParams struct{}

func (JWKSParams) CacheKey() string {
	return cachekeys.JWKSCacheKey()
}
