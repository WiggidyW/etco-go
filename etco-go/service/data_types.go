package service

import (
	"context"

	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"
	"github.com/WiggidyW/etco-go/staticdb"
)

func (Service) AllTypes(
	_ context.Context,
	_ *proto.EmptyRequest,
) (
	rep *proto.TypesResponse,
	_ error,
) {
	rep = &proto.TypesResponse{Types: staticdb.GetSDETypeIds()}
	return rep, nil
}

func (Service) AllNamedTypes(
	_ context.Context,
	_ *proto.EmptyRequest,
) (
	rep *proto.NamedTypesResponse,
	_ error,
) {
	typeDatas := staticdb.UNSAFE_GetSDETypeDatas()
	namedTypes := make([]*proto.NamedTypeId, len(typeDatas))
	r := protoregistry.NewProtoRegistry(len(typeDatas))

	i := 0
	for typeId, typeData := range typeDatas {
		namedTypes[i] = r.UNSAFE_AddType(typeId, typeData)
		i++
	}

	rep = &proto.NamedTypesResponse{Types: namedTypes, Strs: r.Finish()}
	return rep, nil
}

func (Service) NamedTypes(
	_ context.Context,
	req *proto.NamedTypesRequest,
) (
	rep *proto.NamedTypesResponse,
	_ error,
) {
	typeDatas := staticdb.UNSAFE_GetSDETypeDatas()
	namedTypes := make([]*proto.NamedTypeId, len(req.Types))
	r := protoregistry.NewProtoRegistry(len(req.Types))

	for i, typeId := range req.Types {
		if typeData, ok := typeDatas[typeId]; ok {
			namedTypes[i] = r.UNSAFE_AddType(typeId, typeData)
		} else {
			namedTypes[i] = r.AddUndefinedType(typeId)
		}
	}

	rep = &proto.NamedTypesResponse{Types: namedTypes, Strs: r.Finish()}
	return rep, nil
}
