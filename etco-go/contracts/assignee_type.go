package contracts

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
