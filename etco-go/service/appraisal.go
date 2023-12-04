package service

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/appraisal"
	"github.com/WiggidyW/etco-go/auth"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
	"github.com/WiggidyW/etco-go/protoregistry"
)

// TODO: Appraisal can be created before authorized, but not saved
func (Service) SaveBuybackAppraisal(
	ctx context.Context,
	req *proto.SaveAppraisalRequest,
) (
	rep *proto.BuybackAppraisalResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.BuybackAppraisalResponse{}

	var characterId *int32
	rep.Authorized, characterId, err = isSaveAppraisalAuthorized(
		x,
		req.RefreshToken,
	)
	if !rep.Authorized || err != nil {
		rep.Error = protoerr.ErrToProto(err)
		rep.Strs = r.Finish()
		return rep, nil
	}

	rep.Appraisal, _, err = appraisal.ProtoCreateBuybackAppraisal(
		x,
		r,
		req.Items,
		characterId,
		int32(req.TerritoryId),
		true,
	)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
	}

	rep.Strs = r.Finish()
	return rep, nil
}

// TODO: Appraisal can be created before authorized, but not saved
func (Service) SaveShopAppraisal(
	ctx context.Context,
	req *proto.SaveAppraisalRequest,
) (
	rep *proto.ShopAppraisalResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.ShopAppraisalResponse{}

	var characterId *int32
	rep.Authorized, characterId, err = isSaveAppraisalAuthorized(
		x,
		req.RefreshToken,
	)
	if !rep.Authorized || err != nil {
		rep.Error = protoerr.ErrToProto(err)
		rep.Strs = r.Finish()
		return rep, nil
	}

	rep.Appraisal, rep.Status, _, err = appraisal.ProtoCreateShopAppraisal(
		x,
		r,
		req.Items,
		characterId,
		req.TerritoryId,
		true,
	)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
	}

	rep.Strs = r.Finish()
	return rep, nil
}

// TODO: Appraisal can be created before authorized, but not saved
func (Service) SaveHaulAppraisal(
	ctx context.Context,
	req *proto.SaveHaulAppraisalRequest,
) (
	rep *proto.HaulAppraisalResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.HaulAppraisalResponse{}

	var characterId *int32
	rep.Authorized, characterId, err = isSaveAppraisalAuthorized(
		x,
		req.RefreshToken,
	)
	if !rep.Authorized || err != nil {
		rep.Error = protoerr.ErrToProto(err)
		rep.Strs = r.Finish()
		return rep, nil
	}

	rep.Appraisal, _, err = appraisal.ProtoCreateHaulAppraisal(
		x,
		r,
		req.Items,
		characterId,
		req.StartSystemId, req.EndSystemId,
		true,
		req.FallbackPrice,
	)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
	}

	rep.Strs = r.Finish()
	return rep, nil
}

func (Service) NewBuybackAppraisal(
	ctx context.Context,
	req *proto.NewAppraisalRequest,
) (
	rep *proto.BuybackAppraisalResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.BuybackAppraisalResponse{Authorized: true}

	rep.Appraisal, _, err = appraisal.ProtoCreateBuybackAppraisal(
		x,
		r,
		req.Items,
		nil,
		int32(req.TerritoryId),
		false,
	)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
	}

	rep.Strs = r.Finish()
	return rep, nil
}

func (Service) NewShopAppraisal(
	ctx context.Context,
	req *proto.NewAppraisalRequest,
) (
	rep *proto.ShopAppraisalResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.ShopAppraisalResponse{Authorized: true}

	rep.Appraisal, rep.Status, _, err = appraisal.ProtoCreateShopAppraisal(
		x,
		r,
		req.Items,
		nil,
		req.TerritoryId,
		false,
	)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
	}

	rep.Strs = r.Finish()
	return rep, nil
}

func (Service) NewHaulAppraisal(
	ctx context.Context,
	req *proto.NewHaulAppraisalRequest,
) (
	rep *proto.HaulAppraisalResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.HaulAppraisalResponse{Authorized: true}

	rep.Appraisal, _, err = appraisal.ProtoCreateHaulAppraisal(
		x,
		r,
		req.Items,
		nil,
		req.StartSystemId, req.EndSystemId,
		false,
		nil,
	)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
	}

	rep.Strs = r.Finish()
	return rep, nil
}

func (Service) GetBuybackAppraisal(
	ctx context.Context,
	req *proto.GetAppraisalRequest,
) (
	rep *proto.GetBuybackAppraisalResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.GetBuybackAppraisalResponse{}

	rep.Appraisal, rep.Anonymous, err =
		getAppraisal(x, req, r, appraisal.ProtoGetBuybackAppraisal)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
	}

	rep.Strs = r.Finish()
	return rep, nil
}

func (Service) GetShopAppraisal(
	ctx context.Context,
	req *proto.GetAppraisalRequest,
) (
	rep *proto.GetShopAppraisalResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.GetShopAppraisalResponse{}

	rep.Appraisal, rep.Anonymous, err =
		getAppraisal(x, req, r, appraisal.ProtoGetShopAppraisal)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
	}

	rep.Strs = r.Finish()
	return rep, nil
}

func (Service) GetHaulAppraisal(
	ctx context.Context,
	req *proto.GetAppraisalRequest,
) (
	rep *proto.GetHaulAppraisalResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.GetHaulAppraisalResponse{}

	rep.Appraisal, rep.Anonymous, err =
		getAppraisal(x, req, r, appraisal.ProtoGetHaulAppraisal)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
	}

	rep.Strs = r.Finish()
	return rep, nil
}

func isSaveAppraisalAuthorized(
	x cache.Context,
	refreshToken string,
) (
	authorized bool,
	characterId *int32,
	err error,
) {
	if refreshToken == "" {
		return true, nil, nil
	}
	var authRep auth.AuthResponse
	authRep, _, err = auth.ProtoUserAuthorized(x, refreshToken)
	authorized = authRep.Authorized
	characterId = &authRep.CharacterId
	return authorized, characterId, err
}

func getAppraisal[A proto.Appraisal](
	x cache.Context,
	req *proto.GetAppraisalRequest,
	registry *protoregistry.ProtoRegistry,
	get func(
		x cache.Context,
		registry *protoregistry.ProtoRegistry,
		code string,
		include_items bool,
	) (
		appraisal A,
		expires time.Time,
		err error,
	),
) (
	appraisal A,
	anonymous bool,
	err error,
) {
	// fetch characterId+isAdmin in a goroutine
	x, cancel := x.WithCancel()
	defer cancel()
	chnAuth := expirable.NewChanResult[auth.AuthResponse](x.Ctx(), 1, 0)
	go expirable.P2Transceive(
		chnAuth,
		x, req.RefreshToken,
		auth.ProtoAdminAuthorized,
	)

	// fetch appraisal
	appraisal, _, err = get(x, registry, req.Code, req.IncludeItems)
	if appraisal.IsNil() || err != nil || appraisal.GetCharacterId() == 0 {
		// don't check authorization if NilAppraisal, Error, or Anonymous
		return appraisal, true, err
	}

	// recv characterId+isAdmin
	var authRp auth.AuthResponse
	authRp, _, err = chnAuth.RecvExp()
	if !authRp.Authorized && authRp.CharacterId != appraisal.GetCharacterId() {
		// censor characterId if not admin + not same character
		appraisal.ClearCharacterId()
	}
	return appraisal, false, err
}
