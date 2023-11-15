package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CfgMergeShopLocations(
	ctx context.Context,
	req *proto.CfgMergeShopLocationsRequest,
) (rep *proto.CfgMergeShopLocationsResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CfgMergeShopLocationsResponse{}

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

	mergeRep, err := s.cfgMergeShopLocationsClient.Fetch(
		x,
		protoclient.CfgMergeShopLocationsParams{Updates: req.Locations},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
	} else if mergeRep.MergeError != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_INVALID_MERGE,
			mergeRep.MergeError.Error(),
		)
	} else {
		rep.Modified = mergeRep.Modified
	}

	return rep, nil
}
