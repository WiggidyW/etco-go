package esi

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"

	"github.com/lestrrat-go/jwx/jwk"
)

const (
	JWKS_BUF_CAP int = 667
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
