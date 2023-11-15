package proto

import (
	"github.com/WiggidyW/etco-go/staticdb"
)

type AppraisalWithCharacter[A any] struct {
	Appraisal   *A
	CharacterId int32
}

func (awc AppraisalWithCharacter[A]) Unwrap() (
	appraisal *A,
	characterId int32,
) {
	return awc.Appraisal, awc.CharacterId
}

type PBGetAppraisalParams[IM staticdb.IndexMap] struct {
	TypeNamingSession *staticdb.TypeNamingSession[IM]
	AppraisalCode     string
}
