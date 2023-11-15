package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/purchasequeue"
)

func (s *Service) ShopCancelPurchase(
	ctx context.Context,
	req *proto.ShopCancelPurchaseRequest,
) (
	rep *proto.ShopCancelPurchaseResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.ShopCancelPurchaseResponse{}

	var ok bool
	var characterId int32
	characterId, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
		x,
		req.Auth,
		"user",
		true,
	)
	if !ok {
		return rep, nil
	}

	rCancelPurchaseRep, err := purchasequeue.UserCancelPurchase(
		x,
		characterId,
		req.Code,
		0, // TODO
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.Status = protoutil.NewPBCancelPurchaseStatus(rCancelPurchaseRep)

	return rep, nil
}
