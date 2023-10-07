package service

import (
	"context"

	"github.com/WiggidyW/etco-go/client/esi/model/allianceinfo"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) AllianceInfo(
	ctx context.Context,
	req *proto.AllianceInfoRequest,
) (
	rep *proto.AllianceInfoResponse,
	err error,
) {
	rep = &proto.AllianceInfoResponse{}

	rRep, err := s.rAllianceInfoClient.Fetch(
		ctx,
		allianceinfo.AllianceInfoParams{
			AllianceId: req.AllianceId,
		},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.AllianceId = req.AllianceId
	rep.Name = rRep.Data().Name
	rep.Ticker = rRep.Data().Ticker
	return rep, nil
}
