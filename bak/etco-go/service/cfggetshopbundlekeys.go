package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CfgGetShopBundleKeys(
	ctx context.Context,
	req *proto.CfgGetShopBundleKeysRequest,
) (
	rep *proto.CfgGetShopBundleKeysResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CfgGetShopBundleKeysResponse{}

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

	bundleKeysPtr, err := s.cfgGetShopBundleKeysClient.Fetch(
		x,
		protoclient.CfgGetShopBundleKeysParams{},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.BundleKeys = *bundleKeysPtr

	return rep, nil
}
