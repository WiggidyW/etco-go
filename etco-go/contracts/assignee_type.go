package contracts

import (
	"github.com/WiggidyW/etco-go/proto"
)

type AssigneeType uint8

const (
	UnknownAssigneeType AssigneeType = iota
	Corporation
	Character
	Alliance
)

func atFromString(s string) AssigneeType {
	switch s {
	case "corporation":
		return Corporation
	case "personal":
		return Character
	case "alliance":
		return Alliance
	default:
		return UnknownAssigneeType
	}
}

func (at AssigneeType) ToProto() proto.ContractAssigneeType {
	switch at {
	case UnknownAssigneeType:
		return proto.ContractAssigneeType_CAT_UNKNOWN
	case Corporation:
		return proto.ContractAssigneeType_CAT_CORPORATION
	case Character:
		return proto.ContractAssigneeType_CAT_CHARACTER
	case Alliance:
		return proto.ContractAssigneeType_CAT_ALLIANCE
	default:
		panic("Unknown AssigneeType")
	}
}
