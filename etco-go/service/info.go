package service

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
)

func entityInfo[I any](
	x cache.Context,
	req *proto.EntityInfoRequest,
	getInfo func(cache.Context, int32) (*I, time.Time, error),
	notFoundMsg string,
) (
	info *I,
	errResponse *proto.ErrorResponse,
) {
	var err error
	info, _, err = getInfo(x, req.EntityId)
	if err != nil {
		errResponse = protoerr.ErrToProto(err)
	} else if info == nil {
		errResponse = protoerr.MsgNew(protoerr.NOT_FOUND, notFoundMsg).ToProto()
	}
	return info, errResponse
}

func (Service) CharacterInfo(
	ctx context.Context,
	req *proto.EntityInfoRequest,
) (
	rep *proto.CharacterInfoResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CharacterInfoResponse{}

	const notFoundMsg = "Character not found"
	var info *esi.CharacterInfo
	info, rep.Error = entityInfo(x, req, esi.GetCharacterInfo, notFoundMsg)
	if rep.Error != nil {
		return rep, nil
	}

	rep.CharacterId = req.EntityId
	rep.CorporationId = info.CorporationId
	rep.Name = info.Name
	if info.AllianceId != nil {
		rep.AllianceId = *info.AllianceId
	}
	return rep, nil
}

func (Service) CorporationInfo(
	ctx context.Context,
	req *proto.EntityInfoRequest,
) (
	rep *proto.CorporationInfoResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CorporationInfoResponse{}

	const notFoundMsg = "Corporation not found"
	var info *esi.CorporationInfo
	info, rep.Error = entityInfo(x, req, esi.GetCorporationInfo, notFoundMsg)
	if rep.Error != nil {
		return rep, nil
	}

	rep.CorporationId = req.EntityId
	rep.Name = info.Name
	rep.Ticker = info.Ticker
	if info.AllianceId != nil {
		rep.AllianceId = *info.AllianceId
	}
	return rep, nil
}

func (Service) AllianceInfo(
	ctx context.Context,
	req *proto.EntityInfoRequest,
) (
	rep *proto.AllianceInfoResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.AllianceInfoResponse{}

	const notFoundMsg = "Alliance not found"
	var info *esi.AllianceInfo
	info, rep.Error = entityInfo(x, req, esi.GetAllianceInfo, notFoundMsg)
	if rep.Error != nil {
		return rep, nil
	}

	rep.AllianceId = req.EntityId
	rep.Name = info.Name
	rep.Ticker = info.Ticker
	return rep, nil
}
