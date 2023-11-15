package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/remotedb"
)

func (s *Service) UserData(
	ctx context.Context,
	req *proto.UserDataRequest,
) (
	rep *proto.UserDataResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.UserDataResponse{}

	var characterId int32
	var ok bool
	characterId, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
		x,
		req.Auth,
		"user",
		true,
	)
	if !ok {
		return rep, nil
	}

	var rUserData remotedb.UserData
	rUserData, _, err = remotedb.GetUserData(x, characterId)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.BuybackAppraisals,
		rep.ShopAppraisals,
		rep.CancelledPurchase,
		rep.MadePurchase = protoutil.NewPBUserData(rUserData)

	return rep, nil
}
