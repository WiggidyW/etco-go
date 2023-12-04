package service

import (
	"context"

	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"
	"github.com/WiggidyW/etco-go/staticdb"

	b "github.com/WiggidyW/etco-go-bucket"
)

func (Service) AllHaulRoutes(
	_ context.Context,
	_ *proto.EmptyRequest,
) (
	rep *proto.HaulRoutesResponse,
	_ error,
) {
	systemInfos := staticdb.UNSAFE_GetSystemInfos()
	haulRoutes := staticdb.UNSAFE_GetHaulRoutes()
	routes := make([]*proto.HaulRouteInfo, len(haulRoutes))
	r := protoregistry.NewProtoRegistry(len(haulRoutes))

	i := 0
	for key := range haulRoutes {
		startSystemId, endSystemId := staticdb.GetHaulRouteSystems(key)
		routes[i] = &proto.HaulRouteInfo{
			StartSystemInfo: registryAddSystem(systemInfos, r, startSystemId),
			EndSystemInfo:   registryAddSystem(systemInfos, r, endSystemId),
		}
		i++
	}

	rep = &proto.HaulRoutesResponse{Routes: routes, Strs: r.Finish()}
	return rep, nil
}

func (Service) AllSystems(
	_ context.Context,
	_ *proto.EmptyRequest,
) (
	rep *proto.SystemsResponse,
	_ error,
) {
	systemInfos := staticdb.UNSAFE_GetSystemInfos()
	systems := make([]*proto.SystemInfo, len(systemInfos))
	r := protoregistry.NewProtoRegistry(len(systemInfos))

	i := 0
	for systemId, systemInfo := range systemInfos {
		systems[i] = r.UNSAFE_AddSystem(systemId, systemInfo)
		i++
	}

	rep = &proto.SystemsResponse{Systems: systems, Strs: r.Finish()}
	return rep, nil
}

func (Service) AllBuybackSystems(
	_ context.Context,
	_ *proto.EmptyRequest,
) (
	rep *proto.SystemsResponse,
	_ error,
) {
	return systems(staticdb.GetBuybackSystemIds()), nil
}

func (Service) Systems(
	_ context.Context,
	req *proto.SystemsRequest,
) (
	rep *proto.SystemsResponse,
	_ error,
) {
	return systems(req.Systems), nil
}

func systems(systemIds []int32) *proto.SystemsResponse {
	systemInfos := staticdb.UNSAFE_GetSystemInfos()
	systems := make([]*proto.SystemInfo, len(systemIds))
	r := protoregistry.NewProtoRegistry(len(systemIds))
	for i, systemId := range systemIds {
		systems[i] = registryAddSystem(systemInfos, r, systemId)
	}
	return &proto.SystemsResponse{Systems: systems, Strs: r.Finish()}
}

func registryAddSystem(
	systemInfos map[int32]b.System,
	r *protoregistry.ProtoRegistry,
	systemId int32,
) *proto.SystemInfo {
	if systemInfo, ok := systemInfos[systemId]; ok {
		return r.UNSAFE_AddSystem(systemId, systemInfo)
	} else {
		return r.AddUndefinedSystem(systemId)
	}
}
