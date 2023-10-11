package service

import (
	"context"

	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) ShopInventory(
	ctx context.Context,
	req *proto.ShopInventoryRequest,
) (
	rep *proto.ShopInventoryResponse,
	err error,
) {
	rep = &proto.ShopInventoryResponse{}

	var ok bool
	_, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
		ctx,
		req.Auth,
		"user",
		true,
	)
	if !ok {
		return rep, nil
	}

	typeNamingSession := protoutil.
		MaybeNewSyncTypeNamingSession(req.IncludeTypeNaming)

	rep.Items, err = s.shopInventoryClient.Fetch(
		ctx,
		protoclient.PBShopInventoryParams{
			TypeNamingSession: typeNamingSession,
			LocationId:        req.LocationId,
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
