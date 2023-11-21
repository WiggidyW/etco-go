package service

import (
	"context"

	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"
	"github.com/WiggidyW/etco-go/staticdb"
)

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
		if systemInfo, ok := systemInfos[systemId]; ok {
			systems[i] = r.UNSAFE_AddSystem(systemId, systemInfo)
		} else {
			systems[i] = r.AddUndefinedSystem(systemId)
		}
	}

	return &proto.SystemsResponse{Systems: systems, Strs: r.Finish()}
}
