package proto

import (
	"context"

	"github.com/WiggidyW/etco-go/client/appraisal"
	"github.com/WiggidyW/etco-go/proto"
	pu "github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBNewShopAppraisalClient[IM staticdb.IndexMap] struct {
	rNewShopAppraisalClient appraisal.MakeShopAppraisalClient
}

func (nbac PBNewShopAppraisalClient[IM]) Fetch(
	ctx context.Context,
	params PBNewShopAppraisalParams[IM],
) (*proto.ShopAppraisal, error) {
	rAppraisal, err := nbac.rNewShopAppraisalClient.Fetch(
		ctx,
		appraisal.MakeShopAppraisalParams{
			Items:       params.Items,
			LocationId:  params.LocationId,
			CharacterId: params.CharacterId,
			IncludeCode: params.IncludeCode,
		},
	)
	if err != nil {
		return nil, err
	} else {
		return pu.NewPBShopAppraisal(
			*rAppraisal,
			params.TypeNamingSession,
		), nil
	}
}
