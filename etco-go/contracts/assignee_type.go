package contracts

import (
	"github.com/WiggidyW/etco-go/proto"
)

type EntityKind uint8

const (
	EKUnknown EntityKind = iota
	EKCharacter
	EKCorporation
	EKAlliance
)

func entityKindFromStr(s string) EntityKind {
	switch s {
	case "personal":
		return EKCharacter
	case "corporation":
		return EKCorporation
	case "alliance":
		return EKAlliance
	default:
		return EKUnknown
	}
}

func (ek EntityKind) ToProto() proto.EntityKind {
	switch ek {
	case EKUnknown:
		return proto.EntityKind_EK_UNKNOWN
	case EKCharacter:
		return proto.EntityKind_EK_CHARACTER
	case EKCorporation:
		return proto.EntityKind_EK_CORPORATION
	case EKAlliance:
		return proto.EntityKind_EK_ALLIANCE
	default:
		panic("Unknown AssigneeType")
	}
}
