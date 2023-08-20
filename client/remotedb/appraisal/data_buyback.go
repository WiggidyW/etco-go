package appraisal

import "time"

const (
	BUYBACK_COLLECTION_ID string = "buyback_appraisals"

	B_APPR_ITEMS        string = "items"
	B_APPR_PRICE        string = "price"
	B_APPR_TIME         string = "time"
	B_APPR_VERSION      string = "version"
	B_APPR_LOCATION_ID  string = "location_id"
	B_APPR_CHARACTER_ID string = "character_id"
)

type BuybackAppraisal struct {
	Items       []BuybackParentItem `firestore:"items"`
	Price       float64             `firestore:"price"`
	Time        time.Time           `firestore:"time"`
	Version     string              `firestore:"version"`
	LocationId  int64               `firestore:"location_id"`
	CharacterId *int32              `firestore:"character_id"`
}

type BuybackParentItem struct {
	TypeId       int32              `firestore:"type_id"`
	Quantity     int64              `firestore:"quantity"`
	PricePerUnit float64            `firestore:"price_per_unit"`
	Description  string             `firestore:"description"`
	Children     []BuybackChildItem `firestore:"children"`
}

type BuybackChildItem struct {
	TypeId       int32   `firestore:"type_id"`
	Quantity     float64 `firestore:"quantity"`
	PricePerUnit float64 `firestore:"price_per_unit"`
	Description  string  `firestore:"description"`
}

type IBuybackAppraisal[I any] interface {
	IAppraisal[I]
	GetSystemId() int32
}

type IBuybackParentItem[C any] interface {
	IAppraisalItem
	GetQuantity() int64
	GetChildren() []C
}

type IBuybackChildItem interface {
	IAppraisalItem
	GetQuantity() float64
}

func NewBuybackParentItems[
	I IBuybackParentItem[CI],
	CI IBuybackChildItem,
](iParentItems []I) []BuybackParentItem {
	dbParentItems := make([]BuybackParentItem, 0, len(iParentItems))

	for _, iParentItem := range iParentItems {

		iChildItems := iParentItem.GetChildren()
		dbChildItems := make([]BuybackChildItem, 0, len(iChildItems))

		for _, iChildItem := range iChildItems {
			dbChildItems = append(dbChildItems, BuybackChildItem{
				TypeId:       iChildItem.GetTypeId(),
				Quantity:     iChildItem.GetQuantity(),
				PricePerUnit: iChildItem.GetPricePerUnit(),
				Description:  iChildItem.GetDescription(),
			})
		}

		dbParentItems = append(dbParentItems, BuybackParentItem{
			TypeId:       iParentItem.GetTypeId(),
			Quantity:     iParentItem.GetQuantity(),
			PricePerUnit: iParentItem.GetPricePerUnit(),
			Description:  iParentItem.GetDescription(),
			Children:     dbChildItems,
		})
	}

	return dbParentItems
}
