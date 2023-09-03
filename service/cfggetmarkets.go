package service

import (
	"context"

	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) CfgGetMarkets(
	ctx context.Context,
	req *proto.CfgGetMarketsRequest,
) (rep *proto.CfgGetMarketsResponse,
	err error,
) {
	rep = &proto.CfgGetMarketsResponse{}

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

	locationInfoSession := protoutil.MaybeNewSyncLocationInfoSession(
		req.IncludeLocationInfo,
		req.IncludeLocationNaming,
	)

	partialRep, err := s.cfgGetMarketsClient.Fetch(
		ctx,
		protoclient.CfgGetMarketsParams{
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

	rep.Markets = partialRep.Markets
	rep.LocationInfoMap = partialRep.LocationInfoMap
	rep.LocationNamingMaps = protoutil.MaybeFinishLocationInfoSession(
		locationInfoSession,
	)

	return rep, nil
}
