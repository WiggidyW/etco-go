package service

import (
	"context"

	"github.com/WiggidyW/etco-go/client/esi/model/characterinfo"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CharacterInfo(
	ctx context.Context,
	req *proto.CharacterInfoRequest,
) (
	rep *proto.CharacterInfoResponse,
	err error,
) {
	rep = &proto.CharacterInfoResponse{}

	rRep, err := s.rCharacterInfoClient.Fetch(
		ctx,
		characterinfo.CharacterInfoParams{
			CharacterId: req.CharacterId,
		},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.CharacterId = req.CharacterId
	rep.CorporationId = rRep.Data().CorporationId
	rep.Name = rRep.Data().Name
	if rRep.Data().AllianceId != nil {
		rep.AllianceId = &proto.OptionalInt32{
			Inner: *rRep.Data().AllianceId,
		}
	}
	return rep, nil
}
