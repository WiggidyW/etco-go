package service

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/auth"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/contractqueue"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
	"github.com/WiggidyW/etco-go/protoregistry"
	"github.com/WiggidyW/etco-go/purchasequeue"
)

func (Service) CancelPurchase(
	ctx context.Context,
	req *proto.CancelPurchaseRequest,
) (
	rep *proto.CancelPurchaseResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CancelPurchaseResponse{}

	var authRep auth.AuthResponse
	authRep, _, err = auth.ProtoUserAuthorized(x, req.RefreshToken)
	rep.Authorized = authRep.Authorized
	if err != nil || !rep.Authorized {
		rep.Error = protoerr.ErrToProto(err)
		return rep, nil
	}

	rep.Status, err = purchasequeue.ProtoUserCancelPurchase(
		x,
		authRep.CharacterId,
		req.Code,
		req.LocationId,
	)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
	}
	return rep, nil
}

func (Service) DeletePurchases(
	ctx context.Context,
	req *proto.DeletePurchasesRequest,
) (
	rep *proto.DeletePurchasesResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.DeletePurchasesResponse{}

	rep.Authorized, _, err = auth.ProtoBoolAdminAuthorized(x, req.RefreshToken)
	if err != nil || !rep.Authorized {
		rep.Error = protoerr.ErrToProto(err)
		return rep, nil
	}

	err = purchasequeue.DelPurchases(x, req.Entries...)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
	}
	return rep, nil
}

func (Service) LocationPurchaseQueue(
	ctx context.Context,
	req *proto.LocationPurchaseQueueRequest,
) (
	rep *proto.LocationPurchaseQueueResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.LocationPurchaseQueueResponse{}

	rep.Authorized, rep.Queue, err = authorizedGetP2(
		x,
		req.RefreshToken,
		auth.ProtoBoolAdminAuthorized,
		purchasequeue.ProtoGetLocationPurchaseQueue,
		r,
		req.LocationId,
	)
	if !rep.Authorized || err != nil {
		rep.Error = protoerr.ErrToProto(err)
	}

	rep.Strs = r.Finish()
	return rep, nil
}

func authorizedGetQueue[Q any](
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	req *proto.BasicRequest,
	getQueue func(
		cache.Context,
		*protoregistry.ProtoRegistry,
	) (Q, time.Time, error),
) (
	authorized bool,
	queue Q,
	errResponse *proto.ErrorResponse,
) {
	var err error
	authorized, queue, err = authorizedGetP1(
		x,
		req.RefreshToken,
		auth.ProtoBoolAdminAuthorized,
		getQueue,
		r,
	)
	if !authorized || err != nil {
		errResponse = protoerr.ErrToProto(err)
	}
	return authorized, queue, errResponse
}

func (Service) PurchaseQueue(
	ctx context.Context,
	req *proto.BasicRequest,
) (
	rep *proto.PurchaseQueueResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.PurchaseQueueResponse{}
	rep.Authorized, rep.Queue, rep.Error = authorizedGetQueue(
		x,
		r,
		req,
		purchasequeue.ProtoGetPurchaseQueue,
	)
	rep.Strs = r.Finish()
	return rep, nil
}

func (Service) BuybackContractQueue(
	ctx context.Context,
	req *proto.BasicRequest,
) (
	rep *proto.BuybackContractQueueResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.BuybackContractQueueResponse{}
	rep.Authorized, rep.Queue, rep.Error = authorizedGetQueue(
		x,
		r,
		req,
		contractqueue.ProtoGetBuybackContractQueue,
	)
	rep.Strs = r.Finish()
	return rep, nil
}

func (Service) ShopContractQueue(
	ctx context.Context,
	req *proto.BasicRequest,
) (
	rep *proto.ShopContractQueueResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.ShopContractQueueResponse{}
	rep.Authorized, rep.Queue, rep.Error = authorizedGetQueue(
		x,
		r,
		req,
		contractqueue.ProtoGetShopContractQueue,
	)
	rep.Strs = r.Finish()
	return rep, nil
}
