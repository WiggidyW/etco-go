package service

import (
	"github.com/WiggidyW/etco-go/authorized"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) TryAuthenticate(
	x cache.Context,
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

	rAuthRep, _, err := authorized.Authorized(x, authReq.Token, authDomain)
	if err != nil {
		errRep = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return 0, nil, nil, nil, errRep, false
	}
	authRep = NewProtoAuthRep(authReq, rAuthRep)

	if authRep.Authorized {
		return rAuthRep.CharacterId,
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
	authReq *proto.AuthRequest,
	rAuthRep authorized.AuthResponse,
) *proto.AuthResponse {
	return &proto.AuthResponse{
		Token:      authReq.Token,
		Authorized: rAuthRep.Authorized,
	}
}
