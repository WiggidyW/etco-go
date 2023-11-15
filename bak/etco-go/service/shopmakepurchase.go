package service

import (
	"context"

	"github.com/WiggidyW/etco-go/appraisal"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/remotedb"
)

func (s *Service) ShopMakePurchase(
	ctx context.Context,
	req *proto.ShopMakePurchaseRequest,
) (
	rep *proto.ShopMakePurchaseResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.ShopMakePurchaseResponse{}

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

	var rAppraisal remotedb.ShopAppraisal
	rAppraisal, _, err = appraisal.CreateShopAppraisal(
		x,
		req.GetItems(),
		&characterId,
		req.LocationId,
		true,
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	var status appraisal.MakePurchaseStatus
	rAppraisal, status, err = appraisal.SaveShopAppraisal(
		x,
		rAppraisal,
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

	rep.Appraisal = protoutil.NewPBShopAppraisal(
		rAppraisal,
		typeNamingSession,
	)
	rep.Status = protoutil.NewPBMakePurchaseStatus(status)
	rep.TypeNamingLists = protoutil.MaybeFinishTypeNamingSession(
		typeNamingSession,
	)

	return rep, nil
}
