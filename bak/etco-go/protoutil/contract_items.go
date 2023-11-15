package protoutil

import (
	"github.com/WiggidyW/etco-go/contractitems"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/staticdb"
)

func NewPBContractItems[T staticdb.IndexMap](
	rItems []contractitems.ContractItem,
	namingSession *staticdb.TypeNamingSession[T],
) []*proto.ContractItem {
	if len(rItems) == 0 {
		return []*proto.ContractItem{}
	}

	pbItems := make([]*proto.ContractItem, 0, len(rItems))

	for _, rItem := range rItems {
		pbItems = append(pbItems, &proto.ContractItem{
			TypeId:   rItem.TypeId,
			Quantity: rItem.Quantity,
			TypeNamingIndexes: MaybeGetTypeNamingIndexes(
				namingSession,
				rItem.TypeId,
			),
		})
	}

	return pbItems
}
