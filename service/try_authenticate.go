package service

import (
	"context"

	"github.com/WiggidyW/etco-go/client/auth"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) TryAuthenticate(
	ctx context.Context,
	authReq *proto.AuthRequest,
	authDomain string,
	useExtraIds bool,
) (
	characterId int32,
	corporationId *int32,
	allianceId *int32,
	authRep *proto.AuthResponse,
	errRep *proto.ErrorResponse,
	ok bool,
) {
	if authReq == nil {
		errRep = NewProtoErrorRep(
			proto.ErrorCode_INVALID_REQUEST,
			"missing authentication parameters",
		)
		return 0, nil, nil, nil, errRep, false
	}

	rAuthRep, err := s.authClient.Fetch(
		ctx,
		auth.AuthParams{
			NativeRefreshToken: authReq.Token,
			AuthDomain:         authDomain,
			UseExtraIds:        useExtraIds,
		},
	)
	if err != nil {
		errRep = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return 0, nil, nil, nil, errRep, false
	}
	authRep = NewProtoAuthRep(rAuthRep)

	if authRep.Authorized {
		return *rAuthRep.CharacterId,
			rAuthRep.CorporationId,
			rAuthRep.AllianceId,
			authRep,
			nil,
			true
	} else {
		return 0, nil, nil, authRep, nil, false
	}
}

func NewProtoErrorRep(
	code proto.ErrorCode,
	msg string,
) *proto.ErrorResponse {
	return &proto.ErrorResponse{
		Error: msg,
		Code:  code,
	}
}

func NewProtoAuthRep(
	rAuthRep auth.AuthResponse,
) *proto.AuthResponse {
	var token string
	if rAuthRep.NativeRefreshToken != nil {
		token = *rAuthRep.NativeRefreshToken
	}
	return &proto.AuthResponse{
		Token:      token,
		Authorized: rAuthRep.Authorized,
	}
}
