package service

import (
	"context"

	"github.com/WiggidyW/etco-go/auth"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
)

func (Service) TokenInfo(
	ctx context.Context,
	req *proto.TokenInfoRequest,
) (
	rep *proto.TokenInfoResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.TokenInfoResponse{}

	var esiApp esi.EsiApp
	esiApp, err = esi.AppFromProto(req.App)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
		return rep, nil
	}

	var tokenInfoRep auth.ProtoTokenInfoRep
	tokenInfoRep, _, err = auth.ProtoTokenInfo(x, esiApp, req.RefreshToken)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
	} else {
		rep.CharacterId = tokenInfoRep.CharacterId
		rep.Admin = tokenInfoRep.Admin
	}
	return rep, nil
}

func (Service) Login(
	ctx context.Context,
	req *proto.LoginRequest,
) (
	rep *proto.LoginResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.LoginResponse{}

	var esiApp esi.EsiApp
	esiApp, err = esi.AppFromProto(req.App)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
		return rep, nil
	}

	var loginRep auth.ProtoLoginRep
	loginRep, _, err = auth.ProtoLogin(x, esiApp, req.AccessCode)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
	} else {
		rep.RefreshToken = loginRep.RefreshToken
		rep.CharacterId = loginRep.CharacterId
		rep.Admin = loginRep.Admin
	}
	return rep, nil
}
