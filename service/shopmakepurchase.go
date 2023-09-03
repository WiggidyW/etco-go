package service

import (
	"context"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/client/appraisal"
	"github.com/WiggidyW/etco-go/client/purchase"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) ShopMakePurchase(
	ctx context.Context,
	req *proto.ShopMakePurchaseRequest,
) (
	rep *proto.ShopMakePurchaseResponse,
	err error,
) {
	rep = &proto.ShopMakePurchaseResponse{}

	var ok bool
	var characterId int32
	characterId, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
		ctx,
		req.Auth,
		"make-purchase",
		true,
	)
	if !ok {
		return rep, nil
	}

	rBasicItems := make([]appraisal.BasicItem, 0, len(req.Items))
	for _, item := range req.Items {
		rBasicItems = append(rBasicItems, appraisal.BasicItem{
			TypeId:   item.TypeId,
			Quantity: item.Quantity,
		})
	}

	rMakePurchaseRep, err := s.shopMakePurchaseClient.Fetch(
		ctx,
		purchase.MakePurchaseParams{
			Items:       rBasicItems,
			LocationId:  req.LocationId,
			CharacterId: characterId,
			Cooldown:    build.MAKE_PURCHASE_COOLDOWN,
			MaxActive:   build.PURCHASE_MAX_ACTIVE,
		},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	typeNamingSession := protoutil.
		MaybeNewLocalTypeNamingSession(req.IncludeTypeNaming)
	var pbAppraisal *proto.ShopAppraisal
	if rMakePurchaseRep.Appraisal != nil {
		pbAppraisal = protoutil.NewPBShopAppraisal(
			*rMakePurchaseRep.Appraisal,
			typeNamingSession,
		)
	}

	rep.Status = protoutil.NewPBMakePurchaseStatus(
		rMakePurchaseRep.Status,
	)
	rep.Appraisal = pbAppraisal
	rep.TypeNamingLists = protoutil.MaybeFinishTypeNamingSession(
		typeNamingSession,
	)

	return rep, nil
}
