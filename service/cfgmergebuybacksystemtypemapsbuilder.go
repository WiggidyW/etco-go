package service

import (
	"context"

	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CfgMergeBuybackSystemTypeMapsBuilder(
	ctx context.Context,
	req *proto.CfgMergeBuybackSystemTypeMapsBuilderRequest,
) (rep *proto.CfgMergeBuybackSystemTypeMapsBuilderResponse,
	err error,
) {
	rep = &proto.CfgMergeBuybackSystemTypeMapsBuilderResponse{}

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

	mergeRep, err := s.cfgMergeBTypeMapsBuilderClient.Fetch(
		ctx,
		protoclient.CfgMergeBuybackSystemTypeMapsBuilderParams{
			Updates: req.Builder,
		},
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
