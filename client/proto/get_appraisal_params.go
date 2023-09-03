package proto

import (
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBGetAppraisalParams[IM staticdb.IndexMap] struct {
	TypeNamingSession *staticdb.TypeNamingSession[IM]
	AppraisalCode     string
}
