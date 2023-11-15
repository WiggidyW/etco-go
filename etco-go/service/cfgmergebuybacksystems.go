package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CfgMergeBuybackSystems(
	ctx context.Context,
	req *proto.CfgMergeBuybackSystemsRequest,
) (rep *proto.CfgMergeBuybackSystemsResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CfgMergeBuybackSystemsResponse{}

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

	mergeRep, err := s.cfgMergeBuybackSystemsClient.Fetch(
		x,
		protoclient.CfgMergeBuybackSystemsParams{Updates: req.Systems},
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
