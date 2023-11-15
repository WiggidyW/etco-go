package protoutil

import (
	"github.com/WiggidyW/etco-go/contracts"
	"github.com/WiggidyW/etco-go/proto"
)

func NewPBContract(rContract contracts.Contract) *proto.Contract {
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
