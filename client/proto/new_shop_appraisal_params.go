package proto

import (
	"github.com/WiggidyW/etco-go/client/appraisal"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBNewShopAppraisalParams[IM staticdb.IndexMap] struct {
	TypeNamingSession *staticdb.TypeNamingSession[IM]
	Items             []appraisal.BasicItem
	LocationId        int64
	CharacterId       int32
	IncludeCode       bool
}
