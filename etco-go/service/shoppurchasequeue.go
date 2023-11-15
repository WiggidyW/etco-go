package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) ShopPurchaseQueue(
	ctx context.Context,
	req *proto.ShopPurchaseQueueRequest,
) (
	rep *proto.ShopPurchaseQueueResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.ShopPurchaseQueueResponse{}

	var ok bool
	_, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
		x,
		req.Auth,
		"admin",
		true,
	)
	if !ok {
		return rep, nil
	}

	rep.Queue, err = s.shopPurchaseQueueClient.Fetch(
		x,
		protoclient.PBShopPurchaseQueueParams{},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
	}

	return rep, nil
}
