package jwt

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"

	"github.com/WiggidyW/weve-esi/client/esi/internal/raw"
	"github.com/WiggidyW/weve-esi/client/esi/jwt/internal/decodejwks"
	"github.com/WiggidyW/weve-esi/util"
)

// TODO: write a JWKS + JWT library that doesn't use
// excessive reflection

type JWTClient struct {
	jwksClient decodejwks.DecodeJWKSClient
	rawClient  raw.RawClient
}

func (jwtc JWTClient) Fetch(
	ctx context.Context,
	params JWTParams,
) (*JWTResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch JWKS in a goroutine
	chnSend, chnRecv := util.NewChanResult[jwk.Set](ctx).Split()
	go jwtc.fetchJWKS(ctx, chnSend)

	// fetch ESI authentication rep with refresh
	authRep, err := jwtc.rawClient.FetchAuthWithRefresh(
		ctx,
		params.NativeRefreshToken,
	)
	if err != nil {
		return nil, err
	}

	// from hereon, we should return the new native token with errors
	jwtRep := &JWTResponse{NativeRefreshToken: authRep.RefreshToken}

	// wait for JWKS
	jwks, err := chnRecv.Recv()
	if err != nil {
		return jwtRep, err
	}

	// parse the jwt
	jwtRep.CharacterId, err = jwtc.parseJWT(ctx, authRep.AccessToken, jwks)
	return jwtRep, err
}

func (jwtc JWTClient) parseJWT(
	ctx context.Context,
	rawToken string, // JWT token
	jwks jwk.Set,
) (*int32, error) {
	parsedToken, err := jwt.ParseWithClaims(
		rawToken,
		&jWTClaims{},
		func(t *jwt.Token) (interface{}, error) {
			// get the kid
			iKid, ok := t.Header["kid"]
			if !ok { // doesn't exist
				return nil, fmt.Errorf(
					"jwt: kid header not present",
				)
			}
			kid, ok := iKid.(string)
			if !ok { // not a string
				return nil, fmt.Errorf(
					"jwt: kid header not a string",
				)
			}

			// get the jwk
			jwk, ok := jwks.LookupKeyID(kid)
			if !ok { // doesn't exist
				return nil, fmt.Errorf(
					"jwt: jwk not found for kid %s",
					kid,
				)
			}
			var jwtKey interface{}
			if err := jwk.Raw(&jwtKey); err != nil { // invalid jwk
				return nil, err
			}

			return jwtKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &parsedToken.Claims.(*jWTClaims).CharacterID, nil
}

func (jwtc JWTClient) fetchJWKS(
	ctx context.Context,
	chnSend util.ChanSendResult[jwk.Set],
) {
	jwks, err := jwtc.jwksClient.Fetch(ctx, struct{}{})
	if err != nil {
		chnSend.SendErr(err)
	} else {
		chnSend.SendOk(jwks)
	}
}
