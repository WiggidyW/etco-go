package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) BuybackContractQueue(
	ctx context.Context,
	req *proto.BuybackContractQueueRequest,
) (
	rep *proto.BuybackContractQueueResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.BuybackContractQueueResponse{}

	var ok bool
	_, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
		x,
		req.Auth,
		"user",
		true,
	)
	if !ok {
		return rep, nil
	}

	locationInfoSession := protoutil.MaybeNewSyncLocationInfoSession(
		req.IncludeLocationInfo,
		req.IncludeLocationNaming,
	)

	rep.Queue, err = s.buybackContractQueueClient.Fetch(
		x,
		protoclient.PBBuybackContractQueueParams{
			LocationInfoSession: locationInfoSession,
		},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.LocationNamingMaps = protoutil.MaybeFinishLocationInfoSession(
		locationInfoSession,
	)

	return rep, nil
}
