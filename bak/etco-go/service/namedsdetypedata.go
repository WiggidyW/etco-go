package service

import (
	"context"

	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) NamedSDETypeData(
	ctx context.Context,
	req *proto.NamedSDETypeDataRequest,
) (
	rep *proto.NamedSDETypeDataResponse,
	err error,
) {
	rep = &proto.NamedSDETypeDataResponse{}
	rep.Types, rep.TypeNamingLists = protoutil.NewSDENamedTypes()
	return rep, nil
}
