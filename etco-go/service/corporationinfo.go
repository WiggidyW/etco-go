package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CorporationInfo(
	ctx context.Context,
	req *proto.CorporationInfoRequest,
) (
	rep *proto.CorporationInfoResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CorporationInfoResponse{}

	rRep, _, err := esi.GetCorporationInfo(x, req.CorporationId)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.CorporationId = req.CorporationId
	rep.Name = rRep.Name
	rep.Ticker = rRep.Ticker
	if rRep.AllianceId != nil {
		rep.AllianceId = &proto.OptionalInt32{
			Inner: *rRep.AllianceId,
		}
	}
	return rep, nil
}
