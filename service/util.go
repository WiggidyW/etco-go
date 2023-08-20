package service

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/authingfwding"
	"github.com/WiggidyW/weve-esi/client/contracts"
	"github.com/WiggidyW/weve-esi/client/contracts/items"
	"github.com/WiggidyW/weve-esi/client/contracts/items/multi"
	"github.com/WiggidyW/weve-esi/proto"
	"github.com/WiggidyW/weve-esi/staticdb"
)

func authRepToGrpcRep[T any](
	rep *authingfwding.AuthingRep[T],
	err error,
) (ok bool, auth *proto.AuthResponse, e *proto.ErrorResponse) {
	if rep == nil {
		if err == nil {
			panic("unreachable")
		}
		// error = err, authorized = false
		return false, nil, newErrorResponse(err)
	} else if err != nil {
		// error = err, authorized = ?
		return false, newAuthResponse[T](rep), newErrorResponse(err)
	} else if !rep.Authorized {
		// error = nil, authorized = false
		return false, newAuthResponse[T](rep), nil
	} else {
		// error = nil, authorized = true
		return true, newAuthResponse[T](rep), nil
	}
}

func newAuthResponse[T any](
	authRep *authingfwding.AuthingRep[T],
) *proto.AuthResponse {
	return &proto.AuthResponse{
		Authorized: authRep.Authorized,
		Token:      authRep.NativeRefreshToken,
	}
}

func newErrorResponse(err error) *proto.ErrorResponse {
	return &proto.ErrorResponse{
		Error: err.Error(),
		Code:  0, // TODO
	}
}

func maybeNewLocalNamingSession(
	includeNaming *proto.IncludeTypeNaming,
) *staticdb.NamingSession[*staticdb.LocalIndexMap] {
	if includeNaming == nil {
		return nil
	}
	namingSessionVal := staticdb.NewLocalNamingSession(
		includeNaming.IncludeName,
		includeNaming.IncludeMarketGroups,
		includeNaming.IncludeGroup,
		includeNaming.IncludeCategory,
	)
	return &namingSessionVal
}

func maybeNewSyncNamingSession(
	includeNaming *proto.IncludeTypeNaming,
) *staticdb.NamingSession[*staticdb.SyncIndexMap] {
	if includeNaming == nil {
		return nil
	}
	namingSessionVal := staticdb.NewSyncNamingSession(
		includeNaming.IncludeName,
		includeNaming.IncludeMarketGroups,
		includeNaming.IncludeGroup,
		includeNaming.IncludeCategory,
	)
	return &namingSessionVal
}

// func maybeFinishLocalNamingSession(
// 	namingSession *staticdb.NamingSession[*staticdb.LocalIndexMap],
// ) *proto.TypeNamingLists {
// 	return maybeFinishNamingSession(namingSession)
// }

// func maybeFinishSyncNamingSession(
// 	namingSession *staticdb.NamingSession[*staticdb.SyncIndexMap],
// ) *proto.TypeNamingLists {
// 	return maybeFinishNamingSession(namingSession)
// }

func maybeFinishNamingSession[T staticdb.IndexMap](
	namingSession *staticdb.NamingSession[T],
) *proto.TypeNamingLists {
	if namingSession == nil {
		return nil
	}
	mrktGroups, groups, categories := namingSession.Finish()
	return &proto.TypeNamingLists{
		MarketGroups: mrktGroups,
		Groups:       groups,
		Categories:   categories,
	}
}

func newPBTypeNamingIndexes(rNaming staticdb.Naming) *proto.TypeNamingIndexes {
	return &proto.TypeNamingIndexes{
		Name:               rNaming.Name,
		GroupIndex:         rNaming.GroupIndex,
		CategoryIndex:      rNaming.CategoryIndex,
		MarketGroupIndexes: rNaming.MrktGroupIndexes,
	}
}

func maybeTypeNaming[T staticdb.IndexMap](
	namingSession *staticdb.NamingSession[T],
	typeId int32,
) *proto.TypeNamingIndexes {
	if namingSession == nil {
		return nil
	}
	rNaming := namingSession.AddType(typeId)
	return newPBTypeNamingIndexes(rNaming)
}

func newPBContract(rContract contracts.Contract) *proto.Contract {
	return &proto.Contract{
		Status:       newPBContractStatus(rContract.Status),
		Issued:       rContract.Issued.Unix(),
		Expires:      rContract.Expires.Unix(),
		LocationId:   rContract.LocationId,
		Price:        rContract.Price,
		HasReward:    rContract.HasReward,
		IssuerCorpId: rContract.IssuerCorpId,
		IssuerCharId: rContract.IssuerCharId,
		AssigneeId:   rContract.AssigneeId,
		AssigneeType: newPBAssigneeType(rContract.AssigneeType),
		// Items:       nil,
	}
}

func newPBContractItems(
	rItems []items.ContractItem,
	namingSession *staticdb.NamingSession[*staticdb.LocalIndexMap],
) []*proto.ContractItem {
	if rItems == nil {
		return nil
	} else if len(rItems) == 0 {
		return []*proto.ContractItem{}
	}

	// nameItems := namingSession != nil
	pbItems := make([]*proto.ContractItem, 0, len(rItems))

	for _, rItem := range rItems {
		var pbNaming *proto.TypeNamingIndexes = nil
		if namingSession != nil {
			rNaming := namingSession.AddType(rItem.TypeId)
			pbNaming = newPBTypeNamingIndexes(rNaming)
		}
		pbItems = append(pbItems, &proto.ContractItem{
			TypeId:   rItem.TypeId,
			Quantity: rItem.Quantity,
			Naming:   pbNaming,
		})
	}

	return pbItems
}

func newPBContractStatus(
	rStatus contracts.Status,
) proto.ContractStatus {
	switch rStatus {
	case contracts.UnknownStatus:
		return proto.ContractStatus_unknown_status
	case contracts.Outstanding:
		return proto.ContractStatus_outstanding
	case contracts.InProgress:
		return proto.ContractStatus_in_progress
	case contracts.FinishedIssuer:
		return proto.ContractStatus_finished_issuer
	case contracts.FinishedContractor:
		return proto.ContractStatus_finished_contractor
	case contracts.Finished:
		return proto.ContractStatus_finished
	case contracts.Cancelled:
		return proto.ContractStatus_cancelled
	case contracts.Rejected:
		return proto.ContractStatus_rejected
	case contracts.Failed:
		return proto.ContractStatus_failed
	case contracts.Deleted:
		return proto.ContractStatus_deleted
	case contracts.Reversed:
		return proto.ContractStatus_reversed
	default:
		return proto.ContractStatus_unknown_status
	}
}

func newPBAssigneeType(
	rAssigneeType contracts.AssigneeType,
) proto.ContractAssigneeType {
	switch rAssigneeType {
	case contracts.UnknownAssigneeType:
		return proto.ContractAssigneeType_unknown_assignee_type
	case contracts.Corporation:
		return proto.ContractAssigneeType_corporation
	case contracts.Character:
		return proto.ContractAssigneeType_character
	case contracts.Alliance:
		return proto.ContractAssigneeType_alliance
	default:
		return proto.ContractAssigneeType_unknown_assignee_type
	}
}

func (s *Service) fetchAllContractItems(
	ctx context.Context,
	contractIds []int32,
	localNamingSession *staticdb.NamingSession[*staticdb.LocalIndexMap],
) (map[int32][]*proto.ContractItem, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnRecv, err := s.multiContractItemsClient.Fetch(
		ctx,
		multi.MultiRateLimitingContractItemsParams{
			ContractIds: contractIds,
		},
	)
	if err != nil {
		return nil, err
	}

	contractItems := make(map[int32][]*proto.ContractItem, len(contractIds))

	for i := 0; i < len(contractIds); i++ {
		rep, err := chnRecv.Recv()
		if err != nil {
			return nil, err
		}

		contractItems[rep.ContractId] = newPBContractItems(
			rep.ContractItems,
			localNamingSession,
		)
	}

	return contractItems, nil
}
