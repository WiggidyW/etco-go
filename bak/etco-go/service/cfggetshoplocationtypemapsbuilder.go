package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CfgGetShopLocationTypeMapsBuilder(
	ctx context.Context,
	req *proto.CfgGetShopLocationTypeMapsBuilderRequest,
) (rep *proto.CfgGetShopLocationTypeMapsBuilderResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CfgGetShopLocationTypeMapsBuilderResponse{}

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

	rep.Builder, err = s.cfgGetSTypeMapsBuilderClient.Fetch(
		x,
		protoclient.CfgGetShopLocationTypeMapsBuilderParams{},
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
