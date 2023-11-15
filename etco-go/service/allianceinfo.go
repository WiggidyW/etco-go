package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) AllianceInfo(
	ctx context.Context,
	req *proto.AllianceInfoRequest,
) (
	rep *proto.AllianceInfoResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.AllianceInfoResponse{}

	rRep, _, err := esi.GetAllianceInfo(x, req.AllianceId)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.AllianceId = req.AllianceId
	rep.Name = rRep.Name
	rep.Ticker = rRep.Ticker
	return rep, nil
}
