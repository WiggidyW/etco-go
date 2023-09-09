package decodejwks

import (
	"context"
	"encoding/json"

	"github.com/lestrrat-go/jwx/jwk"

	"github.com/WiggidyW/etco-go/client/esi/jwt/jwks"
)

type DecodeJWKSClient struct {
	Inner jwks.WC_JWKSClient
}

func (djwks DecodeJWKSClient) Fetch(
	ctx context.Context,
	params struct{},
) (jwk.Set, error) {
	// fetch the raw JWKS from cache or ESI
	jwksRep, err := djwks.Inner.Fetch(ctx, jwks.JWKSParams{})
	if err != nil {
		return nil, err
	}

	// unmarshal it into a jwk.Set
	jwksSet := jwk.NewSet()
	if err := json.Unmarshal(jwksRep.Data(), &jwksSet); err != nil {
		return nil, err
	}

	return jwksSet, nil
}
