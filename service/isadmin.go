package service

import (
	"context"

	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) IsAdmin(
	ctx context.Context,
	req *proto.IsAdminRequest,
) (
	rep *proto.IsAdminResponse,
	err error,
) {
	rep = &proto.IsAdminResponse{}

	var ok bool
	_, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
		ctx,
		req.Auth,
		"admin",
		false,
	)

	if ok {
		rep.IsAdmin = true
	} else {
		rep.IsAdmin = false
	}

	return rep, nil
}
