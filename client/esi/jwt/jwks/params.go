package jwks

import "github.com/WiggidyW/etco-go/client/cachekeys"

type JWKSParams struct{}

func (JWKSParams) CacheKey() string {
	return cachekeys.JWKSCacheKey()
}
