package jwks

import "github.com/WiggidyW/weve-esi/client/cachekeys"

type JWKSParams struct{}

func (JWKSParams) CacheKey() string {
	return cachekeys.JWKSCacheKey()
}
