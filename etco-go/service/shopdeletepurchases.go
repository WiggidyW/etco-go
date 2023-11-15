package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
)

func (s *Service) ShopDeletePurchases(
	ctx context.Context,
	req *proto.ShopDeletePurchasesRequest,
) (
	rep *proto.ShopDeletePurchasesResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.ShopDeletePurchasesResponse{}

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

	if _, err = s.shopDeletePurchasesClient.Fetch(
		ctx,
		rdbc.DelPurchasesParams{AppraisalCodes: req.Codes},
	); err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		// return rep, nil
	}

	return rep, nil
}
