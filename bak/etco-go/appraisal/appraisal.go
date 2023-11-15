package appraisal

import (
	"github.com/WiggidyW/etco-go/items"
)

type TerritoryInfo interface {
	GetTaxRate() float64
	GetFeePerM3() float64
}

type AppraisalItem interface {
	items.IBasicItem
	GetPricePerUnit() float64
	GetDescription() string
	GetFeePerUnit() float64
	GetChildrenLength() int
}
