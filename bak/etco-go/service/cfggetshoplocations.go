package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) CfgGetShopLocations(
	ctx context.Context,
	req *proto.CfgGetShopLocationsRequest,
) (rep *proto.CfgGetShopLocationsResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CfgGetShopLocationsResponse{}

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

	locationInfoSession := protoutil.MaybeNewSyncLocationInfoSession(
		req.IncludeLocationInfo,
		req.IncludeLocationNaming,
	)

	partialRep, err := s.cfgGetShopLocationsClient.Fetch(
		x,
		protoclient.CfgGetShopLocationsParams{
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

	rep.Locations = partialRep.Locations
	rep.LocationInfoMap = partialRep.LocationInfoMap
	rep.LocationNamingMaps = protoutil.MaybeFinishLocationInfoSession(
		locationInfoSession,
	)

	return rep, nil
}
