package proto

import (
	"context"

	"github.com/WiggidyW/etco-go/client/appraisal"
	"github.com/WiggidyW/etco-go/proto"
	pu "github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBNewBuybackAppraisalClient[IM staticdb.IndexMap] struct {
	rNewBuybackAppraisalClient appraisal.MakeBuybackAppraisalClient
}

func (nbac PBNewBuybackAppraisalClient[IM]) Fetch(
	ctx context.Context,
	params PBNewBuybackAppraisalParams[IM],
) (*proto.BuybackAppraisal, error) {
	rAppraisal, err := nbac.rNewBuybackAppraisalClient.Fetch(
		ctx,
		appraisal.MakeBuybackAppraisalParams{
			Items:       params.Items,
			SystemId:    params.SystemId,
			CharacterId: params.CharacterId,
			Save:        params.Save,
		},
	)
	if err != nil {
		return nil, err
	} else {
		return pu.NewPBBuybackAppraisal(
			*rAppraisal,
			params.TypeNamingSession,
		), nil
	}
}
