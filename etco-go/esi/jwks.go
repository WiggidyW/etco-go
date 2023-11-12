package esi

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"

	"github.com/lestrrat-go/jwx/jwk"
)

const (
	JWKS_BUF_CAP          int           = 667
	JWKS_LOCK_TTL         time.Duration = 30 * time.Second
	JWKS_LOCK_MAX_BACKOFF time.Duration = 10 * time.Second
)

func init() {
	keys.TypeStrJWKS = cache.RegisterType[[]byte]("jwks", JWKS_BUF_CAP)
}

func GetJWKS(x cache.Context) (
	rep jwk.Set, // this type is not cacheable, so we ^ register []byte instead
	expires time.Time,
	err error,
) {
	return jwksGet(x)
}
