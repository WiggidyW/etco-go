package service

import (
	"context"

	"github.com/WiggidyW/etco-go/client/bucket"
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
	rep = &proto.CfgGetConstDataResponse{}

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

	rConstData, err := s.rReadConstDataClient.Fetch(
		ctx,
		bucket.ConstDataReaderParams{},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.ConstData = protoutil.NewPBConstData(rConstData.Data())

	return rep, nil
}
