package proto

import (
	"context"

	"github.com/WiggidyW/etco-go/client/appraisal"
	"github.com/WiggidyW/etco-go/proto"
	pu "github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBNewBuybackAppraisalParams[IM staticdb.IndexMap] struct {
	TypeNamingSession *staticdb.TypeNamingSession[IM]
	Items             []appraisal.BasicItem
	SystemId          int32
	CharacterId       *int32
	Save              bool
}

type PBNewBuybackAppraisalClient[IM staticdb.IndexMap] struct {
	rNewBuybackAppraisalClient appraisal.MakeBuybackAppraisalClient
}

func NewPBNewBuybackAppraisalClient[IM staticdb.IndexMap](
	rNewBuybackAppraisalClient appraisal.MakeBuybackAppraisalClient,
) PBNewBuybackAppraisalClient[IM] {
	return PBNewBuybackAppraisalClient[IM]{rNewBuybackAppraisalClient}
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
