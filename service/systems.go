package service

import (
	"context"

	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) SDESystems(
	ctx context.Context,
	req *proto.SDESystemsRequest,
) (
	rep *proto.SDESystemsResponse,
	err error,
) {
	rep = &proto.SDESystemsResponse{}
	locationInfoSession := protoutil.MaybeNewLocalLocationInfoSession(
		true,
		req.IncludeLocationNaming,
	)
	rep.Systems = protoutil.NewPBSDESystems(locationInfoSession)
	rep.LocationNamingMaps = protoutil.MaybeFinishLocationInfoSession(
		locationInfoSession,
	)
	return rep, nil
}
