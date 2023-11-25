package service

import (
	"context"

	"github.com/WiggidyW/etco-go/parse"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"
)

func (Service) Parse(
	ctx context.Context,
	req *proto.ParseRequest,
) (
	rep *proto.ParseResponse,
	err error,
) {
	r := protoregistry.NewProtoRegistry(0)
	rep = &proto.ParseResponse{}
	parseRep := parse.ProtoParse(r, req.Text)
	rep.KnownItems = parseRep.KnownItems
	rep.UnknownItems = parseRep.UnknownItems
	rep.Strs = r.Finish()
	return rep, nil
}
