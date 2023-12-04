package service

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/appraisal"
	as "github.com/WiggidyW/etco-go/appraisalstatus"
	"github.com/WiggidyW/etco-go/auth"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
	"github.com/WiggidyW/etco-go/protoregistry"
)

func (Service) StatusBuybackAppraisal(
	ctx context.Context,
	req *proto.StatusAppraisalRequest,
) (
	rep *proto.StatusAppraisalResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = authorizedGetAppraisalStatus(
		x,
		r,
		req,
		appraisal.GetBuybackAppraisalCharacterId,
		as.ProtoGetBuybackAppraisalStatus,
	)
	rep.Strs = r.Finish()
	return rep, nil
}

func (Service) StatusShopAppraisal(
	ctx context.Context,
	req *proto.StatusAppraisalRequest,
) (
	rep *proto.StatusAppraisalResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = authorizedGetAppraisalStatus(
		x,
		r,
		req,
		appraisal.GetShopAppraisalCharacterId,
		as.ProtoGetShopAppraisalStatus,
	)
	rep.Strs = r.Finish()
	return rep, nil
}

func (Service) StatusHaulAppraisal(
	ctx context.Context,
	req *proto.StatusAppraisalRequest,
) (
	rep *proto.StatusAppraisalResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = authorizedGetAppraisalStatus(
		x,
		r,
		req,
		appraisal.GetHaulAppraisalCharacterId,
		as.ProtoGetHaulAppraisalStatus,
	)
	rep.Strs = r.Finish()
	return rep, nil
}

func creatorOrAdmin(
	x cache.Context,
	refreshToken string,
	code string,
	getCreator func(cache.Context, string) (*int32, time.Time, error),
) (
	authorized bool,
	err error,
) {
	// fetch the appraisal character in a goroutine
	x, cancel := x.WithCancel()
	defer cancel()
	chnCreator := expirable.NewChanResult[*int32](x.Ctx(), 1, 0)
	go expirable.P2Transceive(chnCreator, x, code, getCreator)

	// check if user is admin + get user character ID
	var authRep auth.AuthResponse
	authRep, _, err = auth.ProtoAdminAuthorized(x, refreshToken)
	if err != nil {
		return false, err
	} else if authRep.Authorized {
		return true, nil
	}

	// recv appraisal character
	var creatorId *int32
	creatorId, _, err = chnCreator.RecvExp()
	if creatorId == nil || err != nil || *creatorId != authRep.CharacterId {
		return false, err
	} else {
		return true, nil
	}
}

func authorizedGetAppraisalStatus(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	req *proto.StatusAppraisalRequest,
	getCreator func(x cache.Context, code string) (*int32, time.Time, error),
	getStatus func(
		x cache.Context,
		r *protoregistry.ProtoRegistry,
		code string,
		include_items bool,
	) (as.ProtoAppraisalStatusRep, time.Time, error),
) (
	rep *proto.StatusAppraisalResponse,
) {
	rep = &proto.StatusAppraisalResponse{}
	var err error

	// fetch the appraisal status in a goroutine
	x, cancel := x.WithCancel()
	defer cancel()
	chnStatus :=
		expirable.NewChanResult[as.ProtoAppraisalStatusRep](x.Ctx(), 1, 0)
	go expirable.P4Transceive(
		chnStatus,
		x, r, req.Code, req.IncludeItems,
		getStatus,
	)

	// check if user is authorized
	rep.Authorized, err =
		creatorOrAdmin(x, req.RefreshToken, req.Code, getCreator)
	if !rep.Authorized || err != nil {
		rep.Error = protoerr.ErrToProto(err)
		return rep
	}

	// recv the appraisal status
	var statusRep as.ProtoAppraisalStatusRep
	statusRep, _, err = chnStatus.RecvExp()
	if err == nil {
		rep.Status = statusRep.Status
		rep.Contract = statusRep.Contract
		rep.ContractItems = statusRep.ContractItems
	}
	return rep
}
