package service

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/client/authingfwding"
	admin "github.com/WiggidyW/eve-trading-co-go/client/configure/authlist"
	"github.com/WiggidyW/eve-trading-co-go/proto"
)

func (s *Service) CfgGetAuthList(
	ctx context.Context,
	req *proto.CfgGetAuthListRequest,
) (*proto.CfgGetAuthListResponse, error) {
	authListRep, err := s.getAuthListClient.Fetch(
		ctx,
		authingfwding.WithAuthableParams[admin.AdminReadParams]{
			NativeRefreshToken: req.Auth.Token,
			Params: admin.AdminReadParams{
				Domain: req.DomainKey,
			},
		},
	)

	ok, authRep, errRep := authRepToGrpcRep(authListRep, err)
	grpcRep := &proto.CfgGetAuthListResponse{
		Auth:  authRep,
		Error: errRep,
	}
	if !ok {
		return grpcRep, nil
	}

	grpcRep.AuthList = &proto.AuthList{
		CharacterIds:   authListRep.Data.CharacterIDs,
		CorporationIds: authListRep.Data.CorporationIDs,
		AllianceIds:    authListRep.Data.AllianceIDs,
	}
	panic("unimplemented")
}
