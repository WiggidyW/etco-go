package service

import (
	"context"

	"github.com/WiggidyW/etco-go/client/bucket"
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
	rep = &proto.CfgSetConstDataResponse{}

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

	_, err = s.rWriteConstDataClient.Fetch(
		ctx,
		bucket.ConstDataWriterParams{
			ConstData: protoutil.NewRConstData(req.ConstData),
		},
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
