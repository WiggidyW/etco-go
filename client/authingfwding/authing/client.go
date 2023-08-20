package authing

import (
	"context"

	"github.com/WiggidyW/weve-esi/client"
	a "github.com/WiggidyW/weve-esi/client/authingfwding"
	"github.com/WiggidyW/weve-esi/client/authingfwding/internal"
)

type AuthingClient[
	P a.AuthFwdableParams[F], // the outer params type
	F any, // the inner client's params type
	D any, // the inner client's response type
	C client.Client[F, D], // the inner client type
] struct {
	Client C
	internal.AuthingClientInner
}

func (ac AuthingClient[P, F, D, C]) Fetch(
	ctx context.Context,
	params P,
) (*a.AuthingRep[D], error) {
	type rep = a.AuthingRep[D]

	// fetch the JWT rep and check if the character is authorized
	jwtRep, authorized, err := ac.Authorized(
		ctx,
		params.AuthRefreshToken(),
	)
	if !authorized {
		if err == nil {
			// return authorized false, the new refresh token
			return a.NewAuthingRep[D](
				nil,
				false,
				jwtRep.NativeRefreshToken,
			), nil
		} else if jwtRep != nil {
			// return authorized false, the new refresh token, the error
			return a.NewAuthingRep[D](
				nil,
				false,
				jwtRep.NativeRefreshToken,
			), err
		} else {
			// return the error
			return nil, err
		}
	} // else authorized true

	// fetch the inner client's rep
	cRep, err := ac.Client.Fetch(
		ctx,
		params.ToInnerParams(*jwtRep.CharacterId),
	)
	if err != nil {
		// return the error with the new refresh token + authorized true
		return a.NewAuthingRep[D](
			nil,
			true,
			jwtRep.NativeRefreshToken,
		), err
	} else {
		// return the rep with the new refresh token + authorized true
		return a.NewAuthingRep[D](
			cRep,
			true,
			jwtRep.NativeRefreshToken,
		), nil
	}
}
