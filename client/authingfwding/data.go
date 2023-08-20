package authingfwding

type AuthingRep[D any] struct {
	Data               *D // nil if not authorized / error
	Authorized         bool
	NativeRefreshToken string // new Native ESI refresh token
}

func NewAuthingRep[D any](
	data *D,
	authorized bool,
	refreshToken string,
) *AuthingRep[D] {
	return &AuthingRep[D]{
		Data:               data,
		Authorized:         authorized,
		NativeRefreshToken: refreshToken,
	}
}

type FwdingRep[D any] struct {
	Data               *D     // nil if error
	NativeRefreshToken string // new Native ESI refresh token
}

func NewFwdingRep[D any](
	data *D,
	refreshToken string,
) *FwdingRep[D] {
	return &FwdingRep[D]{
		Data:               data,
		NativeRefreshToken: refreshToken,
	}
}
