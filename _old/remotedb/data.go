package db

type IShopAppraisal[I any] interface {
	IAppraisal[I]
	GetLocationId() int64
}

type IShopAppraisalItem interface {
	IAppraisalItem
	GetQuantity() int64
}

type CharacterCodes struct {
	BuybackCodes []string
	ShopCodes    []string
}
