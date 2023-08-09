package authing

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/naiveclient/rawclient"
)

type JWTCharacterRep struct {
	CharacterID  *int32 // only nil if an error is also returned
	RefreshToken string // the user of this server must receive this token at all costs
}

type JWTCharacterParams struct {
	RefreshToken string
}

type JWTCharacterClient struct {
	jwksClient cachingJWKSClient
	rawClient  *rawclient.RawClient
}

func (jwtc *JWTCharacterClient) Fetch(
	ctx context.Context,
	params JWTCharacterParams,
) (*JWTCharacterRep, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the jwk set in a goroutine
	chnOk := make(chan jwk.Set, 1)
	chnErr := make(chan error, 1)
	go jwtc.fetchJWKSet(ctx, chnOk, chnErr)

	// get the authentication rep with refresh
	authRep, err := jwtc.rawClient.FetchAuthWithRefresh(
		ctx,
		params.RefreshToken,
	)
	if err != nil {
		return nil, err
	}

	// wait for the jwk set
	var jwkSet jwk.Set
	select {
	case err := <-chnErr:
		return &JWTCharacterRep{RefreshToken: authRep.RefreshToken}, err
	case jwkSet = <-chnOk:
	}

	// parse the jwt
	token, err := jwt.ParseWithClaims(
		authRep.AccessToken,
		&jWTClaims{},
		func(t *jwt.Token) (interface{}, error) {
			// get the kid
			iKid, ok := t.Header["kid"]
			if !ok {
				return nil, fmt.Errorf(
					"jwt: kid header not present",
				)
			}
			kid, ok := iKid.(string)
			if !ok {
				return nil, fmt.Errorf(
					"jwt: kid header not a string",
				)
			}
			// get the jwk
			jwk, ok := jwkSet.LookupKeyID(kid)
			if !ok {
				return nil, fmt.Errorf(
					"jwt: jwk not found for kid %s",
					kid,
				)
			}
			var jwtKey interface{}
			if err := jwk.Raw(&jwtKey); err != nil {
				return nil, err
			}
			return jwtKey, nil
		},
	)
	if err != nil {
		return &JWTCharacterRep{RefreshToken: authRep.RefreshToken}, err
	}

	// return the character id and refresh token
	return &JWTCharacterRep{
		CharacterID:  &token.Claims.(*jWTClaims).CharacterID,
		RefreshToken: authRep.RefreshToken,
	}, nil
}

func (jwtc *JWTCharacterClient) fetchJWKSet(
	ctx context.Context,
	chnOk chan<- jwk.Set,
	chnErr chan<- error,
) {
	// get the raw JWKS json bytes
	bytes, err := jwtc.jwksClient.Fetch(ctx, jWKSBytesClientFetchParams{})
	if err != nil {
		chnErr <- err
		return
	}
	// unmarshal them into a jwk set
	jwkSet := jwk.NewSet()
	if err = json.Unmarshal(bytes.Data(), &jwkSet); err != nil {
		chnErr <- err
		return
	}
	// send the jwk set
	chnOk <- jwkSet
}

func NewJWTCharacterIDClient( // TODO: add capacity param and remove bufpool
	rawClient *rawclient.RawClient,
	bufPool *cache.BufferPool,
	clientCache cache.SharedClientCache,
	serverCache cache.SharedServerCache,
	serverLockTTL time.Duration,
	serverLockMaxWait time.Duration,
) *JWTCharacterClient {
	return &JWTCharacterClient{
		rawClient: rawClient,
		jwksClient: client.NewCachingClient[
			jWKSBytesClientFetchParams,
			[]byte,
			cache.ExpirableData[[]byte],
			jWKSBytesClient,
		](
			jWKSBytesClient{rawClient},
			0,
			bufPool,
			clientCache,
			serverCache,
			serverLockTTL,
			serverLockMaxWait,
		),
	}
}

type jWTClaims struct {
	Audience    string `json:"aud,omitempty"`
	Issuer      string `json:"iss,omitempty"`
	Sub         string `json:"sub,omitempty"`
	CharacterID int32  `json:"-"`
}

func (clm *jWTClaims) Valid() error {
	// validate the parsed claims
	if clm.Issuer != "login.eveonline.com" &&
		clm.Issuer != "https://login.eveonline.com" &&
		clm.Issuer != "http://login.eveonline.com" {
		return fmt.Errorf("jwt: invalid issuer")
	}
	if clm.Audience != "EVE Online" {
		return fmt.Errorf("jwt: invalid audience")
	}
	if clm.Sub == "" {
		return fmt.Errorf("jwt: subject missing or empty")
	}

	// extract the characterID from the subject
	var i int = len(clm.Sub) - 1
	for ; i >= 0; i-- {
		if clm.Sub[i] == ':' {
			break
		}

	}
	id64, err := strconv.ParseInt(clm.Sub[i+1:], 10, 32)
	if err != nil /* || id64 > 2147483647 || id64 < 0 */ {
		return fmt.Errorf("jwt: invalid subject")
	}
	clm.CharacterID = int32(id64)

	// return valid
	return nil
}

type cachingJWKSClient = client.CachingClient[
	jWKSBytesClientFetchParams,
	[]byte,
	cache.ExpirableData[[]byte],
	jWKSBytesClient,
]

type jWKSBytesClientFetchParams struct{}

func (jWKSBytesClientFetchParams) CacheKey() string {
	return "JWKSBytes"
}

type jWKSBytesClient struct {
	*rawclient.RawClient
}

// fetch from server
func (jwks jWKSBytesClient) Fetch(
	ctx context.Context,
	params jWKSBytesClientFetchParams,
) (*cache.ExpirableData[[]byte], error) {
	return jwks.FetchJWKSJson(ctx)
}
