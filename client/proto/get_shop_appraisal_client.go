package proto

import (
	"context"

	rdbc "github.com/WiggidyW/etco-go/client/remotedb"
	"github.com/WiggidyW/etco-go/proto"
	pu "github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBGetShopAppraisalClient[IM staticdb.IndexMap] struct {
	rGetShopAppraisalClient rdbc.WC_ReadShopAppraisalClient
}

func (gsac PBGetShopAppraisalClient[IM]) Fetch(
	ctx context.Context,
	params PBGetAppraisalParams[IM],
) (appraisal AppraisalWithCharacter[proto.ShopAppraisal], err error) {
	rAppraisalRep, err := gsac.rGetShopAppraisalClient.Fetch(
		ctx,
		rdbc.ReadAppraisalParams{AppraisalCode: params.AppraisalCode},
	)
	rAppraisal := rAppraisalRep.Data()
	if err != nil {
		return appraisal, err
	} else if rAppraisal == nil {
		return appraisal, nil
	} else {
		return AppraisalWithCharacter[proto.ShopAppraisal]{
			Appraisal: pu.NewPBShopAppraisal(
				*rAppraisal,
				params.TypeNamingSession,
			),
			CharacterId: rAppraisal.CharacterId,
		}, nil
	}
}
