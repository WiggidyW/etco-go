package service

import (
	"context"

	"github.com/WiggidyW/etco-go/client/esi/raw_"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) EsiAppLogin(
	ctx context.Context,
	req *proto.EsiAppLoginRequest,
) (
	rep *proto.EsiAppLoginResponse,
	err error,
) {
	rep = &proto.EsiAppLoginResponse{}

	var ok bool
	_, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
		ctx,
		req.Auth,
		"admin",
		false,
	)
	if !ok {
		return rep, nil
	}

	var rawClient raw_.RawClient
	switch req.App {
	case proto.EsiApp_EA_AUTH:
		rawClient = s.rAuthRawClient
	case proto.EsiApp_EA_CORPORATION:
		rawClient = s.rCorpRawClient
	case proto.EsiApp_EA_STRUCTURE_INFO:
		rawClient = s.rStructureInfoRawClient
	case proto.EsiApp_EA_MARKETS:
		rawClient = s.rMarketsRawClient
	}

	authRep, err := rawClient.FetchAuthWithRefreshFromCode(
		ctx,
		req.Code,
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.Token = authRep.RefreshToken
	rep.Jwt = authRep.AccessToken
	return rep, nil
}
