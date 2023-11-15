package proto

import (
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
	pu "github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBGetShopAppraisalClient[IM staticdb.IndexMap] struct{}

func NewPBGetShopAppraisalClient[IM staticdb.IndexMap]() PBGetShopAppraisalClient[IM] {
	return PBGetShopAppraisalClient[IM]{}
}

func (gsac PBGetShopAppraisalClient[IM]) Fetch(
	x cache.Context,
	params PBGetAppraisalParams[IM],
) (
	appraisal AppraisalWithCharacter[proto.ShopAppraisal],
	err error,
) {
	rAppraisal, _, err := remotedb.GetShopAppraisal(x, params.AppraisalCode)
	if err != nil {
		return appraisal, err
	} else if rAppraisal == nil {
		return appraisal, nil
	} else {
		var characterId int32
		if rAppraisal.CharacterId != nil {
			characterId = *rAppraisal.CharacterId
		} else {
			characterId = 0
		}
		return AppraisalWithCharacter[proto.ShopAppraisal]{
			Appraisal: pu.NewPBShopAppraisal(
				*rAppraisal,
				params.TypeNamingSession,
			),
			CharacterId: characterId,
		}, nil
	}
}
