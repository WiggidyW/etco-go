package service

import (
	"context"

	"github.com/WiggidyW/etco-go/client/bucket"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CfgGetAuthList(
	ctx context.Context,
	req *proto.CfgGetAuthListRequest,
) (
	rep *proto.CfgGetAuthListResponse,
	err error,
) {
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

	rAuthList, err := s.rReadAuthListClient.Fetch(
		ctx,
		bucket.AuthListReaderParams{AuthDomain: req.DomainKey},
	)
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
