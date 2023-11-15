package service

import (
	"context"

	"github.com/WiggidyW/etco-go/parse"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) Parse(
	ctx context.Context,
	req *proto.ParseRequest,
) (
	rep *proto.ParseResponse,
	err error,
) {
	rep = &proto.ParseResponse{}
	rep.KnownItems, rep.UnknownItems = parse.Parse(req.Text)
	return rep, nil
}
