package fwding

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/client"
	a "github.com/WiggidyW/eve-trading-co-go/client/authingfwding"
	"github.com/WiggidyW/eve-trading-co-go/client/esi/jwt"
)

type FwdingClient[
	P a.AuthFwdableParams[F],
	F any,
	D any,
	C client.Client[F, D],
] struct {
	Client C
	Inner  jwt.JWTClient
}

func (fc FwdingClient[P, F, D, C]) Fetch(
	ctx context.Context,
	params P,
) (*a.AuthingRep[D], error) {
	newRep := a.NewFwdingRep[D]

	// fetch the JWT rep
	jwtRep, err := fc.Inner.Fetch(
		ctx,
		jwt.JWTParams{NativeRefreshToken: params.AuthRefreshToken()},
	)
	if err != nil {
		if jwtRep != nil {
			return newRep(nil, jwtRep.NativeRefreshToken), err
		} else {
			return nil, err
		}
	}

	// fetch the inner client's rep
	cRep, err := fc.Client.Fetch(
		ctx,
		params.ToInnerParams(*jwtRep.CharacterId),
	)
	if err != nil {
		return newRep(nil, jwtRep.NativeRefreshToken), err
	} else {
		return newRep(cRep, jwtRep.NativeRefreshToken), nil
	}
}
