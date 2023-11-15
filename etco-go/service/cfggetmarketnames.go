package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) CfgGetMarketNames(
	ctx context.Context,
	req *proto.CfgGetMarketNamesRequest,
) (
	rep *proto.CfgGetMarketNamesResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.CfgGetMarketNamesResponse{}

	var ok bool
	_, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
		x,
		req.Auth,
		"admin",
		false,
	)
	if !ok {
		return rep, nil
	}

	marketNamesPtr, err := s.cfgGetMarketNamesClient.Fetch(
		x,
		protoclient.CfgGetMarketNamesParams{},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.MarketNames = *marketNamesPtr

	return rep, nil
}
