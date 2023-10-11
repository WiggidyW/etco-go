package protoutil

import "github.com/WiggidyW/etco-go/client/appraisal"

type PBBasicItem interface {
	GetTypeId() int32
	GetQuantity() int64
}

func NewRBasicItems[T PBBasicItem](pbItems []T) []appraisal.BasicItem {
	rItems := make([]appraisal.BasicItem, len(pbItems))
	for i, pbItem := range pbItems {
		rItems[i] = appraisal.BasicItem{
			TypeId:   pbItem.GetTypeId(),
			Quantity: pbItem.GetQuantity(),
		}
	}
	return rItems
}
