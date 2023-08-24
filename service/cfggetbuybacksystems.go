package service

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/client/authingfwding"
	cfg "github.com/WiggidyW/eve-trading-co-go/client/configure"
	"github.com/WiggidyW/eve-trading-co-go/proto"
)

func (s *Service) CfgGetBuybackSystems(
	ctx context.Context,
	req *proto.CfgGetBuybackSystemsRequest,
) (*proto.CfgGetBuybackSystemsResponse, error) {
	buybackSystemsRep, err := s.getBuybackSystemsClient.Fetch(
		ctx,
		authingfwding.WithAuthableParams[struct{}]{
			NativeRefreshToken: req.Auth.Token,
		},
	)

	ok, authRep, errRep := authRepToGrpcRep(buybackSystemsRep, err)
	grpcRep := &proto.CfgGetBuybackSystemsResponse{
		Auth:  authRep,
		Error: errRep,
	}
	if !ok {
		return grpcRep, nil
	}

	grpcRep.Systems = newPBBuybackSystems(
		buybackSystemsRep.Data.Data(),
	)
	return grpcRep, nil
}

func newPBBuybackSystems(
	rBuybackSystems cfg.BuybackSystems,
) *proto.BuybackSystems {
	pbBuybackSystems := &proto.BuybackSystems{
		Inner: make(map[int32]*proto.BuybackSystem),
	}
	for k, v := range rBuybackSystems {
		pbBuybackSystems.Inner[k] = &proto.BuybackSystem{
			BundleKey: v.BundleKey,
			M3Fee:     v.M3Fee,
		}
	}
	return pbBuybackSystems
}
