package service

import (
	"context"

	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) SDETypeData(
	ctx context.Context,
	req *proto.SDETypeDataRequest,
) (
	rep *proto.SDETypeDataResponse,
	err error,
) {
	return &proto.SDETypeDataResponse{
		Types: protoutil.NewSDETypes(),
	}, nil
}
