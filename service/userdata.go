package service

import (
	"context"

	"github.com/WiggidyW/etco-go/client/userdata"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) UserData(
	ctx context.Context,
	req *proto.UserDataRequest,
) (
	rep *proto.UserDataResponse,
	err error,
) {
	rep = &proto.UserDataResponse{}

	var characterId int32
	var ok bool
	characterId, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
		ctx,
		req.Auth,
		"user",
		true,
	)
	if !ok {
		return rep, nil
	}

	rUserData, err := s.rUserDataClient.Fetch(
		ctx,
		userdata.UserDataParams{CharacterId: characterId},
	)
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
