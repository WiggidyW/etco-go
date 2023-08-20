package appraisal

import "time"

const (
	SHOP_COLLECTION_ID string = "shop_appraisals"

	S_APPR_ITEMS        string = "items"
	S_APPR_PRICE        string = "price"
	S_APPR_TIME         string = "time"
	S_APPR_VERSION      string = "version"
	S_APPR_LOCATION_ID  string = "location_id"
	S_APPR_CHARACTER_ID string = "character_id"
)

type ShopAppraisal struct {
	Items       []ShopItem `firestore:"items"`
	Price       float64    `firestore:"price"`
	Time        time.Time  `firestore:"time"`
	Version     string     `firestore:"version"`
	LocationId  int64      `firestore:"location_id"`
	CharacterId int32      `firestore:"character_id"`
}

type ShopItem struct {
	TypeId       int32   `firestore:"type_id"`
	Quantity     int64   `firestore:"quantity"`
	PricePerUnit float64 `firestore:"price_per_unit"`
	Description  string  `firestore:"description"`
}

type IShopAppraisal[I any] interface {
	IAppraisal[I]
	GetLocationId() int64
}

type IShopItem interface {
	IAppraisalItem
	GetQuantity() int64
}

func NewShopItems[I IShopItem](iItems []I) []ShopItem {
	dbItems := make([]ShopItem, 0, len(iItems))
	for _, iItem := range iItems {
		dbItems = append(dbItems, ShopItem{
			TypeId:       iItem.GetTypeId(),
			Quantity:     iItem.GetQuantity(),
			PricePerUnit: iItem.GetPricePerUnit(),
			Description:  iItem.GetDescription(),
		})
	}
	return dbItems
}
