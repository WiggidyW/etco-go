package service

import (
	"context"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) CfgSetConstData(
	ctx context.Context,
	req *proto.CfgSetConstDataRequest,
) (
	rep *proto.CfgSetConstDataResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CfgSetConstDataResponse{}

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

	err = bucket.SetBuildConstData(
		x,
		protoutil.NewRConstData(req.ConstData),
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
