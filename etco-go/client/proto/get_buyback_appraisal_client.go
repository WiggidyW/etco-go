package proto

import (
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
	pu "github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBGetBuybackAppraisalClient[IM staticdb.IndexMap] struct{}

func NewPBGetBuybackAppraisalClient[IM staticdb.IndexMap]() PBGetBuybackAppraisalClient[IM] {
	return PBGetBuybackAppraisalClient[IM]{}
}

func (gbac PBGetBuybackAppraisalClient[IM]) Fetch(
	x cache.Context,
	params PBGetAppraisalParams[IM],
) (
	appraisal AppraisalWithCharacter[proto.BuybackAppraisal],
	err error,
) {
	rAppraisal, _, err := remotedb.GetBuybackAppraisal(x, params.AppraisalCode)
	if err != nil {
		return appraisal, err
	} else if rAppraisal == nil { // return nil appraisal
		return appraisal, nil
	} else {
		var characterId int32
		if rAppraisal.CharacterId != nil {
			characterId = *rAppraisal.CharacterId
		} else {
			characterId = 0
		}
		return AppraisalWithCharacter[proto.BuybackAppraisal]{
			Appraisal: pu.NewPBBuybackAppraisal(
				*rAppraisal,
				params.TypeNamingSession,
			),
			CharacterId: characterId,
		}, nil
	}
}
