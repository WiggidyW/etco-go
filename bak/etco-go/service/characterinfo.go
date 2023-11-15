package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CharacterInfo(
	ctx context.Context,
	req *proto.CharacterInfoRequest,
) (
	rep *proto.CharacterInfoResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CharacterInfoResponse{}

	rRep, _, err := esi.GetCharacterInfo(x, req.CharacterId)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.CharacterId = req.CharacterId
	rep.CorporationId = rRep.CorporationId
	rep.Name = rRep.Name
	if rRep.AllianceId != nil {
		rep.AllianceId = &proto.OptionalInt32{
			Inner: *rRep.AllianceId,
		}
	}
	return rep, nil
}
