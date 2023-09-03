package service

import (
	"context"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/client/purchase"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) ShopCancelPurchase(
	ctx context.Context,
	req *proto.ShopCancelPurchaseRequest,
) (
	rep *proto.ShopCancelPurchaseResponse,
	err error,
) {
	rep = &proto.ShopCancelPurchaseResponse{}

	var ok bool
	var characterId int32
	characterId, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
		ctx,
		req.Auth,
		"cancel-purchase",
		true,
	)
	if !ok {
		return rep, nil
	}

	rCancelPurchaseRep, err := s.shopCancelPurchaseClient.Fetch(
		ctx,
		purchase.CancelPurchaseParams{
			AppraisalCode: req.Code,
			CharacterId:   characterId,
			Cooldown:      build.CANCEL_PURCHASE_COOLDOWN,
		},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.Status = protoutil.NewPBCancelPurchaseStatus(*rCancelPurchaseRep)

	return rep, nil
}
