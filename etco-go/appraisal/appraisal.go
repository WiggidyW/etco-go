package appraisal

type TerritoryInfo interface {
	GetTaxRate() float64
	GetFeePerM3() float64
}

type BasicItem interface {
	GetTypeId() int32
	GetQuantity() int64
}

type AppraisalItem interface {
	BasicItem
	GetPricePerUnit() float64
	GetDescription() string
	GetFeePerUnit() float64
	GetChildrenLength() int
}
