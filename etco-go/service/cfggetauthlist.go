package service

import (
	"context"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CfgGetAuthList(
	ctx context.Context,
	req *proto.CfgGetAuthListRequest,
) (
	rep *proto.CfgGetAuthListResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CfgGetAuthListResponse{}

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

	rAuthList, _, err := bucket.GetAuthList(x, req.DomainKey)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.AuthList = &proto.AuthList{
		BannedCharacterIds:   rAuthList.BannedCharacterIds,
		PermitCharacterIds:   rAuthList.PermitCharacterIds,
		BannedCorporationIds: rAuthList.BannedCorporationIds,
		PermitCorporationIds: rAuthList.PermitCorporationIds,
		PermitAllianceIds:    rAuthList.PermitAllianceIds,
	}

	return rep, nil
}
