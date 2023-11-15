package service

import (
	"context"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CfgSetAuthList(
	ctx context.Context,
	req *proto.CfgSetAuthListRequest,
) (
	rep *proto.CfgSetAuthListResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CfgSetAuthListResponse{}

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

	err = bucket.SetAuthList(
		x,
		req.DomainKey,
		bucket.AuthList{
			BannedCharacterIds:   req.AuthList.BannedCharacterIds,
			PermitCharacterIds:   req.AuthList.PermitCharacterIds,
			BannedCorporationIds: req.AuthList.BannedCorporationIds,
			PermitCorporationIds: req.AuthList.PermitCorporationIds,
			PermitAllianceIds:    req.AuthList.PermitAllianceIds,
		},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	return rep, nil
}
