package service

import (
	"context"

	a "github.com/WiggidyW/weve-esi/client/appraisal"
	"github.com/WiggidyW/weve-esi/client/appraisal/shop"
	"github.com/WiggidyW/weve-esi/proto"
	"github.com/WiggidyW/weve-esi/staticdb"
)

func (s *Service) NewShopAppraisal(
	ctx context.Context,
	req *proto.NewShopAppraisalRequest,
) (*proto.NewShopAppraisalResponse, error) {
	apprRep, err := s.admnShopApprClient.Fetch(
		ctx,
		shop.FWD_ShopAppraisalParams{
			NativeRefreshToken: req.Auth.Token,
			Params: shop.INNERFWD_ShopAppraisalParams{
				Items:       newRBasicItems(req.Items),
				LocationId:  req.LocationId,
				IncludeCode: false,
			},
		},
	)
	ok, authRep, errRep := authRepToGrpcRep(apprRep, err)
	grpcRep := &proto.NewShopAppraisalResponse{
		Auth:  authRep,
		Error: errRep,
	}
	if !ok {
		return grpcRep, nil
	}

	namingSession := maybeNewLocalNamingSession(req.IncludeNaming)
	grpcRep.Appraisal = newPBShopAppraisal(*apprRep.Data, namingSession)

	return grpcRep, nil
}

func newPBShopAppraisal[T staticdb.IndexMap](
	rAppraisal a.ShopAppraisal,
	namingSession *staticdb.NamingSession[T],
) *proto.ShopAppraisal {
	pbAppraisal := &proto.ShopAppraisal{
		Items: make(
			[]*proto.ShopItem,
			0,
			len(rAppraisal.Items),
		),
		Code:       "",
		Price:      rAppraisal.Price,
		Time:       rAppraisal.Time.Unix(),
		Version:    rAppraisal.Version,
		LocationId: rAppraisal.LocationId,
		Naming:     nil,
	}
	for _, rItem := range rAppraisal.Items {
		pbAppraisal.Items = append(
			pbAppraisal.Items,
			newPBAppraisalShopItem(rItem, namingSession),
		)
	}
	pbAppraisal.Naming = maybeFinishNamingSession(namingSession)
	return pbAppraisal
}

func newPBAppraisalShopItem[T staticdb.IndexMap](
	rShopItem a.ShopItem,
	namingSession *staticdb.NamingSession[T],
) *proto.ShopItem {
	return &proto.ShopItem{
		TypeId:       rShopItem.TypeId,
		Quantity:     rShopItem.Quantity,
		PricePerUnit: rShopItem.PricePerUnit,
		Description:  rShopItem.Description,
		Naming:       maybeTypeNaming(namingSession, rShopItem.TypeId),
	}
}
