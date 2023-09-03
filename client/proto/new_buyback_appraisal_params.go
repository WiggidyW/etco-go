package proto

import (
	"github.com/WiggidyW/etco-go/client/appraisal"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBNewBuybackAppraisalParams[IM staticdb.IndexMap] struct {
	TypeNamingSession *staticdb.TypeNamingSession[IM]
	Items             []appraisal.BasicItem
	SystemId          int32
	CharacterId       *int32
	Save              bool
}
