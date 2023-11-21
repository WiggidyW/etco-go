package service

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/auth"
	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
)

func cfgGet[CFG any](
	ctx context.Context,
	req *proto.BasicRequest,
	get func(cache.Context) (CFG, time.Time, error),
) (
	authorized bool,
	cfg CFG,
	protoErr *proto.ErrorResponse,
) {
	x := cache.NewContext(ctx)
	var err error
	authorized, cfg, err = authorizedGet(
		x,
		req.RefreshToken,
		auth.ProtoBoolAdminAuthorized,
		get,
	)
	if err != nil {
		protoErr = protoerr.ErrToProto(err)
	}
	return authorized, cfg, protoErr
}

func (Service) CfgGetConstData(
	ctx context.Context,
	req *proto.BasicRequest,
) (
	rep *proto.CfgGetConstDataResponse,
	_ error,
) {
	rep = &proto.CfgGetConstDataResponse{}
	rep.Authorized, rep.ConstData, rep.Error = cfgGet(
		ctx,
		req,
		bucket.ProtoGetBuildConstData,
	)
	return rep, nil
}

func (Service) CfgGetUserAuthList(
	ctx context.Context,
	req *proto.BasicRequest,
) (
	rep *proto.CfgGetAuthListResponse,
	_ error,
) {
	rep = &proto.CfgGetAuthListResponse{}
	rep.Authorized, rep.AuthList, rep.Error = cfgGet(
		ctx,
		req,
		bucket.ProtoGetUserAuthHashSet,
	)
	return rep, nil
}

func (Service) CfgGetAdminAuthList(
	ctx context.Context,
	req *proto.BasicRequest,
) (
	rep *proto.CfgGetAuthListResponse,
	_ error,
) {
	rep = &proto.CfgGetAuthListResponse{}
	rep.Authorized, rep.AuthList, rep.Error = cfgGet(
		ctx,
		req,
		bucket.ProtoGetAdminAuthHashSet,
	)
	return rep, nil
}

func (Service) CfgGetBuybackSystemTypeMapsBuilder(
	ctx context.Context,
	req *proto.BasicRequest,
) (
	rep *proto.CfgGetBuybackSystemTypeMapsBuilderResponse,
	_ error,
) {
	rep = &proto.CfgGetBuybackSystemTypeMapsBuilderResponse{}
	rep.Authorized, rep.Builder, rep.Error = cfgGet(
		ctx,
		req,
		bucket.ProtoGetWebBuybackSystemTypeMapsBuilder,
	)
	return rep, nil
}

func (Service) CfgGetShopLocationTypeMapsBuilder(
	ctx context.Context,
	req *proto.BasicRequest,
) (
	rep *proto.CfgGetShopLocationTypeMapsBuilderResponse,
	_ error,
) {
	rep = &proto.CfgGetShopLocationTypeMapsBuilderResponse{}
	rep.Authorized, rep.Builder, rep.Error = cfgGet(
		ctx,
		req,
		bucket.ProtoGetWebShopLocationTypeMapsBuilder,
	)
	return rep, nil
}

func (Service) CfgGetBuybackSystems(
	ctx context.Context,
	req *proto.BasicRequest,
) (
	rep *proto.CfgGetBuybackSystemsResponse,
	_ error,
) {
	rep = &proto.CfgGetBuybackSystemsResponse{}
	rep.Authorized, rep.Systems, rep.Error = cfgGet(
		ctx,
		req,
		bucket.ProtoGetWebBuybackSystems,
	)
	return rep, nil
}

func (Service) CfgGetShopLocations(
	ctx context.Context,
	req *proto.BasicRequest,
) (
	rep *proto.CfgGetShopLocationsResponse,
	_ error,
) {
	rep = &proto.CfgGetShopLocationsResponse{}
	rep.Authorized, rep.Locations, rep.Error = cfgGet(
		ctx,
		req,
		bucket.ProtoGetWebShopLocations,
	)
	return rep, nil
}

func (Service) CfgGetMarkets(
	ctx context.Context,
	req *proto.BasicRequest,
) (
	rep *proto.CfgGetMarketsResponse,
	_ error,
) {
	rep = &proto.CfgGetMarketsResponse{}
	rep.Authorized, rep.Markets, rep.Error = cfgGet(
		ctx,
		req,
		bucket.ProtoGetWebMarkets,
	)
	return rep, nil
}

func (Service) CfgGetBuybackBundleKeys(
	ctx context.Context,
	req *proto.BasicRequest,
) (
	rep *proto.CfgGetBuybackBundleKeysResponse,
	_ error,
) {
	rep = &proto.CfgGetBuybackBundleKeysResponse{}
	rep.Authorized, rep.BundleKeys, rep.Error = cfgGet(
		ctx,
		req,
		bucket.ProtoGetWebBuybackBundleKeys,
	)
	return rep, nil
}

func (Service) CfgGetShopBundleKeys(
	ctx context.Context,
	req *proto.BasicRequest,
) (
	rep *proto.CfgGetShopBundleKeysResponse,
	_ error,
) {
	rep = &proto.CfgGetShopBundleKeysResponse{}
	rep.Authorized, rep.BundleKeys, rep.Error = cfgGet(
		ctx,
		req,
		bucket.ProtoGetWebShopBundleKeys,
	)
	return rep, nil
}

func (Service) CfgGetMarketNames(
	ctx context.Context,
	req *proto.BasicRequest,
) (
	rep *proto.CfgGetMarketNamesResponse,
	_ error,
) {
	rep = &proto.CfgGetMarketNamesResponse{}
	rep.Authorized, rep.MarketNames, rep.Error = cfgGet(
		ctx,
		req,
		bucket.ProtoGetWebMarketNames,
	)
	return rep, nil
}
