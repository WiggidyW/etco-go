package contracts

import (
	"github.com/WiggidyW/etco-go/proto"
)

type Status uint8

const (
	UnknownStatus Status = iota
	Outstanding
	InProgress
	FinishedIssuer
	FinishedContractor
	Finished
	Cancelled
	Rejected
	Failed
	Deleted
	Reversed
)

func sFromString(s string) Status {
	switch s {
	case "outstanding":
		return Outstanding
	case "in_progress":
		return InProgress
	case "finished_issuer":
		return FinishedIssuer
	case "finished_contractor":
		return FinishedContractor
	case "finished":
		return Finished
	case "cancelled":
		return Cancelled
	case "rejected":
		return Rejected
	case "failed":
		return Failed
	case "deleted":
		return Deleted
	case "reversed":
		return Reversed
	default:
		return UnknownStatus
	}
}

func (s Status) ToProto() proto.ContractStatus {
	switch s {
	case UnknownStatus:
		return proto.ContractStatus_CS_UNKNOWN
	case Outstanding:
		return proto.ContractStatus_CS_OUTSTANDING
	case InProgress:
		return proto.ContractStatus_CS_IN_PROGRESS
	case FinishedIssuer:
		return proto.ContractStatus_CS_FINISHED_ISSUER
	case FinishedContractor:
		return proto.ContractStatus_CS_FINISHED_CONTRACTOR
	case Finished:
		return proto.ContractStatus_CS_FINISHED
	case Cancelled:
		return proto.ContractStatus_CS_CANCELLED
	case Rejected:
		return proto.ContractStatus_CS_REJECTED
	case Failed:
		return proto.ContractStatus_CS_FAILED
	case Deleted:
		return proto.ContractStatus_CS_DELETED
	case Reversed:
		return proto.ContractStatus_CS_REVERSED
	default:
		panic("Unknown Status")
	}
}
