package proto

import (
	"context"

	"github.com/WiggidyW/etco-go/client/contracts"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBContractItemsClient[IM staticdb.IndexMap] struct {
	rSingleContractItemsClient contracts.WC_SingleContractItemsClient
}

func (gcic PBContractItemsClient[IM]) Fetch(
	ctx context.Context,
	params PBContractItemsParams[IM],
) ([]*proto.ContractItem, error) {
	rContractItems, err := gcic.rSingleContractItemsClient.Fetch(
		ctx,
		contracts.SingleContractItemsParams{
			ContractId: params.ContractId,
		},
	)
	if err != nil {
		return nil, err
	} else {
		return protoutil.NewPBContractItems(
			rContractItems.Data(),
			params.TypeNamingSession,
		), nil
	}
}
