package authingfwding

type FwdableParams[F any] interface {
	ToInnerParams(characterId int32) F
}

type AuthFwdableParams[F any] interface {
	FwdableParams[F]
	AuthRefreshToken() string
}

type WithAuthableParams[F any] struct {
	NativeRefreshToken string
	Params             F
}

func (p WithAuthableParams[F]) AuthRefreshToken() string {
	return p.NativeRefreshToken
}

func (p WithAuthableParams[F]) ToInnerParams(characterId int32) F {
	return p.Params
}

type WithAuthFwdableParams[F any, P FwdableParams[F]] struct {
	NativeRefreshToken string
	Params             P
}

func (p WithAuthFwdableParams[P, F]) AuthRefreshToken() string {
	return p.NativeRefreshToken
}

func (p WithAuthFwdableParams[F, P]) ToInnerParams(characterId int32) F {
	return p.Params.ToInnerParams(characterId)
}
