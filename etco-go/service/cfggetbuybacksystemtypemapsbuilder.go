package service

import (
	"context"

	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CfgGetBuybackSystemTypeMapsBuilder(
	ctx context.Context,
	req *proto.CfgGetBuybackSystemTypeMapsBuilderRequest,
) (
	rep *proto.CfgGetBuybackSystemTypeMapsBuilderResponse,
	err error,
) {
	rep = &proto.CfgGetBuybackSystemTypeMapsBuilderResponse{}

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

	rep.Builder, err = s.cfgGetBTypeMapsBuilderClient.Fetch(
		ctx,
		protoclient.CfgGetBuybackSystemTypeMapsBuilderParams{},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	return rep, nil
}
