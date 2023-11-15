package contractitems

import (
	"github.com/WiggidyW/etco-go/esi"
)

type ContractItems = []ContractItem

func fromEntries(entries []esi.ContractItemsEntry) []ContractItem {
	items := make([]ContractItem, 0, len(entries))
	itemsMap := make(map[int32]int64, len(entries))
	for _, entry := range entries {
		itemsMap[entry.TypeId] += int64(entry.Quantity)
	}
	for typeId, quantity := range itemsMap {
		items = append(items, ContractItem{quantity, typeId})
	}
	return items
}
