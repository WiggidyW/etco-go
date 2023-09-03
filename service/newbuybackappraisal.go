package service

import (
	"context"

	"github.com/WiggidyW/etco-go/client/appraisal"
	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

func (s *Service) NewBuybackAppraisal(
	ctx context.Context,
	req *proto.NewBuybackAppraisalRequest,
) (
	rep *proto.NewBuybackAppraisalResponse,
	err error,
) {
	rep = &proto.NewBuybackAppraisalResponse{}

	var characterId *int32
	if req.Auth != nil {
		var characterIdVal int32
		var ok bool
		characterIdVal, _, _, rep.Auth, rep.Error, ok =
			s.TryAuthenticate(
				ctx,
				req.Auth,
				"new-buyback-appraisal",
				true,
			)
		if !ok {
			return rep, nil
		}
		characterId = &characterIdVal
	}

	typeNamingSession := protoutil.
		MaybeNewLocalTypeNamingSession(req.IncludeTypeNaming)
	rep.Appraisal, err = s.newBuybackAppraisalClient.Fetch(
		ctx,
		protoclient.PBNewBuybackAppraisalParams[*staticdb.LocalIndexMap]{
			TypeNamingSession: typeNamingSession,
			Items:             newRBasicItems(req.GetItems()),
			SystemId:          req.SystemId,
			CharacterId:       characterId,
			Save:              req.Save,
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

func newRBasicItems(pbItems []*proto.BasicItem) []appraisal.BasicItem {
	rItems := make([]appraisal.BasicItem, 0, len(pbItems))
	for _, pbItem := range pbItems {
		rItems = append(rItems, appraisal.BasicItem{
			TypeId:   pbItem.TypeId,
			Quantity: pbItem.Quantity,
		})
	}
	return rItems
}
