package service

import (
	"context"

	a "github.com/WiggidyW/weve-esi/client/appraisal"
	"github.com/WiggidyW/weve-esi/client/appraisal/buyback"
	"github.com/WiggidyW/weve-esi/client/authingfwding"
	"github.com/WiggidyW/weve-esi/proto"
	"github.com/WiggidyW/weve-esi/staticdb"
)

func (s *Service) NewBuybackAppraisal(
	ctx context.Context,
	req *proto.NewBuybackAppraisalRequest,
) (*proto.NewBuybackAppraisalResponse, error) {
	grpcRep := &proto.NewBuybackAppraisalResponse{}

	var rAppraisal *a.BuybackAppraisal
	var err error

	if req.Auth != nil {
		var apprRep *authingfwding.AuthingRep[a.BuybackAppraisal]
		apprRep, err = s.charBuybackApprClient.Fetch(
			ctx,
			buyback.FWD_BuybackAppraisalParams{
				NativeRefreshToken: req.Auth.Token,
				Params: buyback.INNERFWD_BuybackAppraisalParams{
					Items:    newRBasicItems(req.Items),
					SystemId: req.SystemId,
					Save:     req.Save,
				},
			},
		)
		if apprRep != nil {
			grpcRep.Auth = newAuthResponse(apprRep)
			if apprRep.Data != nil {
				rAppraisal = apprRep.Data
			}
		}

	} else {
		rAppraisal, err = s.anonBuybackApprClient.Fetch(
			ctx,
			buyback.BuybackAppraisalParams{
				Items:       newRBasicItems(req.Items),
				SystemId:    req.SystemId,
				CharacterId: nil,
				Save:        req.Save,
			},
		)
	}

	if err != nil {
		grpcRep.Error = newErrorResponse(err)
		return grpcRep, nil
	}

	namingSession := maybeNewLocalNamingSession(req.IncludeNaming)
	grpcRep.Appraisal = newPBBuybackAppraisal(*rAppraisal, namingSession)

	return grpcRep, nil
}

func newPBBuybackAppraisal(
	rAppraisal a.BuybackAppraisal,
	namingSession *staticdb.NamingSession[*staticdb.LocalIndexMap],
) *proto.BuybackAppraisal {
	pbAppraisal := &proto.BuybackAppraisal{
		Items: make(
			[]*proto.BuybackParentItem,
			0,
			len(rAppraisal.Items),
		),
		Code:     rAppraisal.Code,
		Price:    rAppraisal.Price,
		Time:     rAppraisal.Time.Unix(),
		Version:  rAppraisal.Version,
		SystemId: rAppraisal.SystemId,
		Naming:   nil,
	}
	for _, rParentItem := range rAppraisal.Items {
		pbAppraisal.Items = append(
			pbAppraisal.Items,
			newPBBuybackParentItem(
				rParentItem,
				namingSession,
			),
		)
	}
	pbAppraisal.Naming = maybeFinishNamingSession(namingSession)
	return pbAppraisal
}

func newPBBuybackParentItem(
	rParentItem a.BuybackParentItem,
	namingSession *staticdb.NamingSession[*staticdb.LocalIndexMap],
) *proto.BuybackParentItem {
	pbParentItem := &proto.BuybackParentItem{
		TypeId:       rParentItem.TypeId,
		Quantity:     rParentItem.Quantity,
		PricePerUnit: rParentItem.PricePerUnit,
		Description:  rParentItem.Description,
		Children: make(
			[]*proto.BuybackChildItem,
			0,
			len(rParentItem.Children),
		),
		Naming: maybeTypeNaming(namingSession, rParentItem.TypeId),
	}
	for _, rChildItem := range rParentItem.Children {
		pbParentItem.Children = append(
			pbParentItem.Children,
			newPBBuybackChildItem(rChildItem, namingSession),
		)
	}
	return pbParentItem
}

func newPBBuybackChildItem(
	rChildItem a.BuybackChildItem,
	namingSession *staticdb.NamingSession[*staticdb.LocalIndexMap],
) *proto.BuybackChildItem {
	return &proto.BuybackChildItem{
		TypeId:            rChildItem.TypeId,
		QuantityPerParent: rChildItem.QuantityPerParent,
		PricePerUnit:      rChildItem.PricePerUnit,
		Description:       rChildItem.Description,
		Naming: maybeTypeNaming(
			namingSession,
			rChildItem.TypeId,
		),
	}
}
