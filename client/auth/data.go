package auth

type AuthResponse struct {
	Authorized         bool
	NativeRefreshToken *string // only nil if returned alongside a non-nil error
	CharacterId        *int32  // only nil if returned alongside a non-nil error
	CorporationId      *int32  // possibly nil if check wasn't needed
	AllianceId         *int32  // possibly nil if check wasn't needed or character not in alliance
}
