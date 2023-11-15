package proto

import (
	"github.com/WiggidyW/etco-go/appraisal"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/items"
	"github.com/WiggidyW/etco-go/proto"
	pu "github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBNewBuybackAppraisalParams[IM staticdb.IndexMap] struct {
	TypeNamingSession *staticdb.TypeNamingSession[IM]
	Items             []items.BasicItem
	SystemId          int32
	CharacterId       *int32
	Save              bool
}

type PBNewBuybackAppraisalClient[IM staticdb.IndexMap] struct{}

func NewPBNewBuybackAppraisalClient[IM staticdb.IndexMap]() PBNewBuybackAppraisalClient[IM] {
	return PBNewBuybackAppraisalClient[IM]{}
}

func (nbac PBNewBuybackAppraisalClient[IM]) Fetch(
	x cache.Context,
	params PBNewBuybackAppraisalParams[IM],
) (*proto.BuybackAppraisal, error) {
	rAppraisal, _, err := appraisal.CreateBuybackAppraisal(
		x,
		params.Items,
		params.CharacterId,
		params.SystemId,
		params.Save,
	)
	if err != nil {
		return nil, err
	}
	if params.Save {
		err = appraisal.SaveBuybackAppraisal(x, rAppraisal)
		if err != nil {
			return nil, err
		}
	}
	return pu.NewPBBuybackAppraisal(
		rAppraisal,
		params.TypeNamingSession,
	), nil
}
