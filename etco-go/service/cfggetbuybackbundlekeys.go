package service

import (
	"context"

	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CfgGetBuybackBundleKeys(
	ctx context.Context,
	req *proto.CfgGetBuybackBundleKeysRequest,
) (
	rep *proto.CfgGetBuybackBundleKeysResponse,
	err error,
) {
	rep = &proto.CfgGetBuybackBundleKeysResponse{}

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

	bundleKeysPtr, err := s.cfgGetBuybackBundleKeysClient.Fetch(
		ctx,
		protoclient.CfgGetBuybackBundleKeysParams{},
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
