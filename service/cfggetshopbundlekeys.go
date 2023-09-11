package service

import (
	"context"

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
	rep = &proto.CfgGetShopBundleKeysResponse{}

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

	bundleKeysPtr, err := s.cfgGetShopBundleKeysClient.Fetch(
		ctx,
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
