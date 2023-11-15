package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	protoclient "github.com/WiggidyW/etco-go/client/proto"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/remotedb"
)

func (s *Service) StatusBuybackAppraisal(
	ctx context.Context,
	req *proto.StatusBuybackAppraisalRequest,
) (
	rep *proto.StatusBuybackAppraisalResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.StatusBuybackAppraisalResponse{}

	var ok bool

	if req.Admin {
		_, _, _, rep.Auth, rep.Error, ok = s.TryAuthenticate(
			x,
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
			x,
			req.Auth,
			"user",
			true,
		)
		if !ok {
			return rep, nil
		}

		rUserCodes, _, err := remotedb.GetUserBuybackAppraisalCodes(
			x,
			characterId,
		)
		if err != nil {
			rep.Error = NewProtoErrorRep(
				proto.ErrorCode_SERVER_ERROR,
				err.Error(),
			)
			return rep, nil
		}

		var userHasCode bool
		for _, code := range rUserCodes {
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
		x,
		protoclient.PBStatusAppraisalParams{
			TypeNamingSession:   typeNamingSession,
			LocationInfoSession: locationInfoSession,
			AppraisalCode:       req.Code,
			IncludeItems:        req.IncludeItems,
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
