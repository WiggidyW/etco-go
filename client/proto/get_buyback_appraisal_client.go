package proto

import (
	"context"

	rdbc "github.com/WiggidyW/etco-go/client/remotedb"
	"github.com/WiggidyW/etco-go/proto"
	pu "github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBGetBuybackAppraisalClient[IM staticdb.IndexMap] struct {
	rGetBuybackAppraisalClient rdbc.WC_ReadBuybackAppraisalClient
}

func NewPBGetBuybackAppraisalClient[IM staticdb.IndexMap](
	rGetBuybackAppraisalClient rdbc.WC_ReadBuybackAppraisalClient,
) PBGetBuybackAppraisalClient[IM] {
	return PBGetBuybackAppraisalClient[IM]{rGetBuybackAppraisalClient}
}

func (gbac PBGetBuybackAppraisalClient[IM]) Fetch(
	ctx context.Context,
	params PBGetAppraisalParams[IM],
) (appraisal AppraisalWithCharacter[proto.BuybackAppraisal], err error) {
	rAppraisalRep, err := gbac.rGetBuybackAppraisalClient.Fetch(
		ctx,
		rdbc.ReadAppraisalParams{AppraisalCode: params.AppraisalCode},
	)
	rAppraisal := rAppraisalRep.Data()
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
