package proto

import (
	"context"

	"github.com/WiggidyW/etco-go/client/contracts"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBContractItemsParams[IM staticdb.IndexMap] struct {
	TypeNamingSession *staticdb.TypeNamingSession[IM]
	ContractId        int32
}

type PBContractItemsClient[IM staticdb.IndexMap] struct {
	rSingleContractItemsClient contracts.WC_SingleContractItemsClient
}

func NewPBContractItemsClient[IM staticdb.IndexMap](
	rSingleContractItemsClient contracts.WC_SingleContractItemsClient,
) PBContractItemsClient[IM] {
	return PBContractItemsClient[IM]{
		rSingleContractItemsClient,
	}
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
