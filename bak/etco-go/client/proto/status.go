package proto

import (
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBStatusAppraisalParams struct {
	TypeNamingSession   *staticdb.TypeNamingSession[*staticdb.LocalIndexMap]
	LocationInfoSession *staticdb.LocationInfoSession[*staticdb.LocalLocationNamerTracker]
	AppraisalCode       string
	IncludeItems        bool
}
