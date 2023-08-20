package appraisal

const (
	CHARACTERS_COLLECTION_ID string = "characters"
)

type IAppraisal[I any] interface {
	GetItems() []I
	GetCode() string
	GetPrice() float64
	GetVersion() string
	GetLocationId() int64
}

type IAppraisalItem interface {
	GetTypeId() int32
	GetPricePerUnit() float64
	GetDescription() string
}
