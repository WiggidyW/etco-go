package items

import (
	"time"

	ci "github.com/WiggidyW/eve-trading-co-go/client/esi/model/contractitems"
)

const (
	LIMITED_CODE        int           = 520
	LIMITED_STR         string        = "ConStopSpamming"
	LIMITED_SLEEP       time.Duration = 10 * time.Second
	MAX_CONCURRENT_REQS int           = 20
	// max: 20 reqs per 10 secs
)

type ContractItem struct {
	TypeId   int32
	Quantity int64
}

// converts entries to items by combining the quantities of same type entries
func EntriesToItems(entries []ci.ContractItemsEntry) []ContractItem {

	// return early if combining is un-necessary
	if len(entries) == 0 {
		return []ContractItem{}
	} else if len(entries) == 1 {
		return []ContractItem{{
			TypeId:   entries[0].TypeId,
			Quantity: int64(entries[0].Quantity),
		}}
	}

	// convert the entries to a map of items for combining
	itemsMap := make(map[int32]*ContractItem, len(entries))
	for _, entry := range entries {
		item, ok := itemsMap[entry.TypeId]
		if !ok {
			itemsMap[entry.TypeId] = &ContractItem{
				TypeId:   entry.TypeId,
				Quantity: int64(entry.Quantity),
			}
		} else {
			item.Quantity += int64(entry.Quantity)
		}
	}

	// then, convert the map to a slice
	itemsSlice := make([]ContractItem, 0, len(itemsMap))
	for _, item := range itemsMap {
		itemsSlice = append(itemsSlice, *item)
	}

	return itemsSlice
}
