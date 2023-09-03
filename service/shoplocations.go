package service

import (
	"context"

	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) ShopLocations(
	ctx context.Context,
	req *proto.ShopLocationsRequest,
) (
	rep *proto.ShopLocationsResponse,
	err error,
) {
	rep = &proto.ShopLocationsResponse{}

	locationInfoSession := protoutil.MaybeNewSyncLocationInfoSession(
		req.IncludeLocationInfo,
		req.IncludeLocationNaming,
	)

	rep.Locations, err = s.shopLocationsClient.Fetch(
		ctx,
		protoclient.ShopLocationsParams{
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

	rep.LocationNamingMaps = protoutil.MaybeFinishLocationInfoSession(
		locationInfoSession,
	)

	return rep, nil
}
