package service

import (
	"context"

	"github.com/WiggidyW/etco-go/client/esi/model/corporationinfo"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CorporationInfo(
	ctx context.Context,
	req *proto.CorporationInfoRequest,
) (
	rep *proto.CorporationInfoResponse,
	err error,
) {
	rep = &proto.CorporationInfoResponse{}

	rRep, err := s.rCorporationInfoClient.Fetch(
		ctx,
		corporationinfo.CorporationInfoParams{
			CorporationId: req.CorporationId,
		},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.CorporationId = req.CorporationId
	rep.Name = rRep.Data().Name
	rep.Ticker = rRep.Data().Ticker
	if rep.AllianceId != nil {
		rep.AllianceId = &proto.OptionalInt32{
			Inner: *rRep.Data().AllianceId,
		}
	}
	return rep, nil
}
