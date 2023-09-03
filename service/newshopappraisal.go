package service

import (
	"context"

	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

func (s *Service) NewShopAppraisal(
	ctx context.Context,
	req *proto.NewShopAppraisalRequest,
) (
	rep *proto.NewShopAppraisalResponse,
	err error,
) {
	rep = &proto.NewShopAppraisalResponse{}

	typeNamingSession := protoutil.
		MaybeNewLocalTypeNamingSession(req.IncludeTypeNaming)
	rep.Appraisal, err = s.newShopAppraisalClient.Fetch(
		ctx,
		protoclient.PBNewShopAppraisalParams[*staticdb.LocalIndexMap]{
			TypeNamingSession: typeNamingSession,
			Items:             newRBasicItems(req.GetItems()),
			LocationId:        req.LocationId,
			CharacterId:       0,
			IncludeCode:       false,
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
