package proto

import (
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/contractitems"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBContractItemsParams[IM staticdb.IndexMap] struct {
	TypeNamingSession *staticdb.TypeNamingSession[IM]
	ContractId        int32
}

type PBContractItemsClient[IM staticdb.IndexMap] struct {
}

func NewPBContractItemsClient[IM staticdb.IndexMap]() PBContractItemsClient[IM] {
	return PBContractItemsClient[IM]{}
}

func (gcic PBContractItemsClient[IM]) Fetch(
	x cache.Context,
	params PBContractItemsParams[IM],
) ([]*proto.ContractItem, error) {
	rContractItems, _, err := contractitems.GetContractItems(
		x,
		params.ContractId,
	)
	if err != nil {
		return nil, err
	} else {
		return protoutil.NewPBContractItems(
			rContractItems,
			params.TypeNamingSession,
		), nil
	}
}
