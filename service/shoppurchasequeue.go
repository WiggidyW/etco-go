package service

import (
	"context"

	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) ShopPurchaseQueue(
	ctx context.Context,
	req *proto.ShopPurchaseQueueRequest,
) (
	rep *proto.ShopPurchaseQueueResponse,
	err error,
) {
	rep = &proto.ShopPurchaseQueueResponse{}

	var ok bool
	_, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
		ctx,
		req.Auth,
		"admin",
		true,
	)
	if !ok {
		return rep, nil
	}

	typeNamingSession := protoutil.MaybeNewSyncTypeNamingSession(
		req.IncludeTypeNaming,
	)

	rep.Queue, err = s.shopPurchaseQueueClient.Fetch(
		ctx,
		protoclient.PBPurchaseQueueParams{
			TypeNamingSession: typeNamingSession,
			QueueInclude: protoclient.NewPurchaseQueueInclude(
				req.IncludeCodeAppraisal,
				req.IncludeNewAppraisal,
			),
		},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.TypeNamingLists = protoutil.MaybeFinishTypeNamingSession(
		typeNamingSession,
	)

	return rep, nil
}
