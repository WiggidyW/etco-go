package service

import (
	"context"

	protoclient "github.com/WiggidyW/etco-go/client/proto"
	rdbc "github.com/WiggidyW/etco-go/client/remotedb"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

func (s *Service) StatusBuybackAppraisal(
	ctx context.Context,
	req *proto.StatusBuybackAppraisalRequest,
) (
	rep *proto.StatusBuybackAppraisalResponse,
	err error,
) {
	rep = &proto.StatusBuybackAppraisalResponse{}

	var ok bool

	if req.Admin {
		_, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
			ctx,
			req.Auth,
			"admin",
			false,
		)
		if !ok {
			return rep, nil
		}

	} else {
		var characterId int32
		characterId, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
			ctx,
			req.Auth,
			"status-buyback-appraisal",
			true,
		)
		if !ok {
			return rep, nil
		}

		rUserDataRep, err := s.rdbcUserDataClient.Fetch(
			ctx,
			rdbc.ReadUserDataParams{CharacterId: characterId},
		)
		if err != nil {
			return rep, err
		}

		var userHasCode bool
		for _, code := range rUserDataRep.Data().BuybackAppraisals {
			if code == req.Code {
				userHasCode = true
				break
			}
		}
		if !userHasCode {
			rep.Auth.Authorized = false
			return rep, nil
		}
	}

	typeNamingSession := protoutil.
		MaybeNewLocalTypeNamingSession(req.IncludeTypeNaming)
	locationInfoSession := protoutil.
		MaybeNewLocalLocationInfoSession(
			req.IncludeLocationInfo,
			req.IncludeLocationNaming,
		)
	statusAppraisalRep, err := s.statusBuybackAppraisalClient.Fetch(
		ctx,
		protoclient.PBStatusAppraisalParams{
			TypeNamingSession:   typeNamingSession,
			LocationInfoSession: locationInfoSession,
			AppraisalCode:       req.Code,
			StatusInclude: protoclient.NewAppraisalStatusInclude(
				req.IncludeItems,
				req.IncludeLocationInfo,
			),
		},
	)
	if err != nil {
		rep.Error = NewProtoErrorRep(
			proto.ErrorCode_SERVER_ERROR,
			err.Error(),
		)
		return rep, nil
	}

	rep.Contract = statusAppraisalRep.Contract
	rep.ContractItems = statusAppraisalRep.ContractItems
	rep.LocationInfo = statusAppraisalRep.LocationInfo
	rep.TypeNamingLists = protoutil.MaybeFinishTypeNamingSession(
		typeNamingSession,
	)
	rep.LocationNamingMaps = protoutil.MaybeFinishLocationInfoSession(
		locationInfoSession,
	)

	return rep, nil
}
