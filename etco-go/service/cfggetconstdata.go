package service

import (
	"context"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) CfgGetConstData(
	ctx context.Context,
	req *proto.CfgGetConstDataRequest,
) (
	rep *proto.CfgGetConstDataResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CfgGetConstDataResponse{}

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

	rConstData, _, err := bucket.GetBuildConstData(x)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.ConstData = protoutil.NewPBConstData(rConstData)

	return rep, nil
}
