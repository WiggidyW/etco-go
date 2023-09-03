package contracts

import (
	"strings"
	"time"

	"github.com/WiggidyW/etco-go/client/appraisal"
	modelci "github.com/WiggidyW/etco-go/client/esi/model/contractitems"
	"github.com/WiggidyW/etco-go/error/esierror"
)

const (
	CI_LIMITED_CODE      int           = 520
	CI_LIMITED_STR       string        = "ConStopSpamming"
	CI_ATTEMPT_INTERVAL  time.Duration = 10 * time.Second
	CI_REQS_PER_INTERVAL int           = 20
	// max: 20 reqs per 10 secs
)

type ContractItem = appraisal.BasicItem

// converts entries to items by combining the quantities of same type entries
func EntriesToItems(entries []modelci.ContractItemsEntry) []ContractItem {

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

func RateLimited(err error) bool {
	statusErr, ok := err.(esierror.StatusError)
	if ok && statusErr.Code == CI_LIMITED_CODE && strings.Contains(
		statusErr.EsiText,
		CI_LIMITED_STR,
	) {
		return true
	}
	return false
}
