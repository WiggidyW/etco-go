package service

import (
	"context"

	"github.com/WiggidyW/etco-go/auth"
	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
)

func cfgSet[CFG any](
	ctx context.Context,
	refreshToken string,
	cfg CFG,
	set func(cache.Context, CFG) error,
) (
	rep *proto.CfgUpdateResponse,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CfgUpdateResponse{}

	var err error
	rep.Authorized, _, err = auth.ProtoBoolAdminAuthorized(x, refreshToken)
	if !rep.Authorized || err != nil {
		rep.Error = protoerr.ErrToProto(err)
		return rep
	}

	err = set(x, cfg)
	if err != nil {
		rep.Error = protoerr.ErrToProto(err)
	} else {
		rep.Modified = true
	}
	return rep
}

func (Service) CfgSetUserAuthList(
	ctx context.Context,
	req *proto.CfgSetAuthListRequest,
) (
	rep *proto.CfgUpdateResponse,
	_ error,
) {
	rep = cfgSet(
		ctx,
		req.RefreshToken,
		req.AuthList,
		bucket.ProtoSetUserAuthHashSet,
	)
	return rep, nil
}

func (Service) CfgSetAdminAuthList(
	ctx context.Context,
	req *proto.CfgSetAuthListRequest,
) (
	rep *proto.CfgUpdateResponse,
	_ error,
) {
	rep = cfgSet(
		ctx,
		req.RefreshToken,
		req.AuthList,
		bucket.ProtoSetAdminAuthHashSet,
	)
	return rep, nil
}

func (Service) CfgSetConstData(
	ctx context.Context,
	req *proto.CfgSetConstDataRequest,
) (
	rep *proto.CfgUpdateResponse,
	_ error,
) {
	rep = cfgSet(
		ctx,
		req.RefreshToken,
		req.ConstData,
		bucket.ProtoSetBuildConstData,
	)
	return rep, nil
}

func (Service) CfgMergeBuybackSystemTypeMapsBuilder(
	ctx context.Context,
	req *proto.CfgMergeBuybackSystemTypeMapsBuilderRequest,
) (
	rep *proto.CfgUpdateResponse,
	_ error,
) {
	rep = cfgSet(
		ctx,
		req.RefreshToken,
		req.Builder,
		bucket.ProtoMergeSetWebBuybackSystemTypeMapsBuilder,
	)
	return rep, nil
}

func (Service) CfgMergeShopLocationTypeMapsBuilder(
	ctx context.Context,
	req *proto.CfgMergeShopLocationTypeMapsBuilderRequest,
) (
	rep *proto.CfgUpdateResponse,
	_ error,
) {
	rep = cfgSet(
		ctx,
		req.RefreshToken,
		req.Builder,
		bucket.ProtoMergeSetWebShopLocationTypeMapsBuilder,
	)
	return rep, nil
}

func (Service) CfgMergeHaulRouteTypeMapsBuilder(
	ctx context.Context,
	req *proto.CfgMergeHaulRouteTypeMapsBuilderRequest,
) (
	rep *proto.CfgUpdateResponse,
	_ error,
) {
	rep = cfgSet(
		ctx,
		req.RefreshToken,
		req.Builder,
		bucket.ProtoMergeSetWebHaulRouteTypeMapsBuilder,
	)
	return rep, nil
}

func (Service) CfgMergeBuybackSystems(
	ctx context.Context,
	req *proto.CfgMergeBuybackSystemsRequest,
) (
	rep *proto.CfgUpdateResponse,
	_ error,
) {
	rep = cfgSet(
		ctx,
		req.RefreshToken,
		req.Systems,
		bucket.ProtoMergeSetWebBuybackSystems,
	)
	return rep, nil
}

func (Service) CfgMergeShopLocations(
	ctx context.Context,
	req *proto.CfgMergeShopLocationsRequest,
) (
	rep *proto.CfgUpdateResponse,
	_ error,
) {
	rep = cfgSet(
		ctx,
		req.RefreshToken,
		req.Locations,
		bucket.ProtoMergeSetWebShopLocations,
	)
	return rep, nil
}

func (Service) CfgMergeHaulRoutes(
	ctx context.Context,
	req *proto.CfgMergeHaulRoutesRequest,
) (
	rep *proto.CfgUpdateResponse,
	_ error,
) {
	rep = cfgSet(
		ctx,
		req.RefreshToken,
		req.Routes,
		bucket.ProtoMergeSetWebHaulRoutes,
	)
	return rep, nil
}

func (Service) CfgMergeMarkets(
	ctx context.Context,
	req *proto.CfgMergeMarketsRequest,
) (
	rep *proto.CfgUpdateResponse,
	_ error,
) {
	rep = cfgSet(
		ctx,
		req.RefreshToken,
		req.Markets,
		bucket.ProtoMergeSetWebMarkets,
	)
	return rep, nil
}
