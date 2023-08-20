package contracts

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
