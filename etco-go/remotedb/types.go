package remotedb

import "time"

type BasicItem struct {
	TypeId   int32
	Quantity int64
}

type ShopQueue struct {
	ShopQueue []string `firestore:"shop_queue"`
}

type UserData struct {
	BuybackAppraisals []string  `firestore:"buyback_appraisals"`
	ShopAppraisals    []string  `firestore:"shop_appraisals"`
	CancelledPurchase time.Time `firestore:"cancelled_purchase"`
	MadePurchase      time.Time `firestore:"made_purchase"`
}

type BuybackAppraisal struct {
	// ignored during reading (used as doc id instead of field)
	// technically, if you're reading, you must already know it
	Code string `firestore:"-"`

	// ignored during writing (we use a nifty serverTimestamp firestore feature instead)
	Time time.Time `firestore:"time"`

	Items       []BuybackParentItem `firestore:"items"`
	Fee         float64             `firestore:"fee,omitempty"`
	Tax         float64             `firestore:"tax,omitempty"`
	Price       float64             `firestore:"price"`
	Version     string              `firestore:"version"`
	SystemId    int32               `firestore:"system_id"`
	CharacterId *int32              `firestore:"character_id"`
}

type BuybackParentItem struct {
	TypeId       int32              `firestore:"type_id"`
	Quantity     int64              `firestore:"quantity"`
	PricePerUnit float64            `firestore:"price_per_unit"`
	FeePerUnit   float64            `firestore:"fee,omitempty"`
	Description  string             `firestore:"description"`
	Children     []BuybackChildItem `firestore:"children"`
}

type BuybackChildItem struct {
	TypeId            int32   `firestore:"type_id"`
	QuantityPerParent float64 `firestore:"quantity_per_parent"`
	PricePerUnit      float64 `firestore:"price_per_unit"`
	Description       string  `firestore:"description"`
}

type ShopAppraisal struct {
	// ignored during reading (used as doc id instead of field)
	// technically, if you're reading, you must already know it
	Code string `firestore:"-"`

	// ignored during writing (we use a nifty serverTimestamp firestore feature instead)
	Time time.Time `firestore:"time"`

	Items       []ShopItem `firestore:"items"`
	Price       float64    `firestore:"price"`
	Tax         float64    `firestore:"tax,omitempty"`
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
