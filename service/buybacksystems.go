package service

import (
	"context"

	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) BuybackSystems(
	ctx context.Context,
	req *proto.BuybackSystemsRequest,
) (
	rep *proto.BuybackSystemsResponse,
	err error,
) {
	rep = &proto.BuybackSystemsResponse{}
	locationInfoSession := protoutil.MaybeNewLocalLocationInfoSession(
		true,
		req.IncludeLocationNaming,
	)
	rep.Systems = protoutil.NewPBBuybackSystems(locationInfoSession)
	rep.LocationNamingMaps = protoutil.MaybeFinishLocationInfoSession(
		locationInfoSession,
	)
	return rep, nil
}
