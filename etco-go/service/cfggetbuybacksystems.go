package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) CfgGetBuybackSystems(
	ctx context.Context,
	req *proto.CfgGetBuybackSystemsRequest,
) (
	rep *proto.CfgGetBuybackSystemsResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CfgGetBuybackSystemsResponse{}

	var ok bool
	_, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
		x,
		req.Auth,
		"admin",
		false,
	)
	if !ok {
		return rep, nil
	}

	locationInfoSession := protoutil.MaybeNewLocalLocationInfoSession(
		req.IncludeLocationInfo,
		req.IncludeLocationNaming,
	)

	partialRep, err := s.cfgGetBuybackSystemsClient.Fetch(
		x,
		protoclient.CfgGetBuybackSystemsParams{
			LocationInfoSession: locationInfoSession,
		},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.Systems = partialRep.Systems
	rep.SystemRegionMap = partialRep.SystemRegionMap
	rep.LocationNamingMaps = protoutil.MaybeFinishLocationInfoSession(
		locationInfoSession,
	)

	return rep, nil
}
