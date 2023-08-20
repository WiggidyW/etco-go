package appraisal

import (
	"time"
)

type BuybackAppraisal struct {
	Items       []BuybackParentItem `firestore:"items"`
	Price       float64             `firestore:"price"`
	Time        time.Time           `firestore:"time"`
	Version     string              `firestore:"version"`
	SystemId    int32               `firestore:"system_id"`
	CharacterId *int32              `firestore:"character_id"`
}

type BuybackParentItem struct {
	TypeId       int32              `firestore:"type_id"`
	Quantity     int64              `firestore:"quantity"`
	PricePerUnit float64            `firestore:"price_per_unit"`
	Fee          float64            `firestore:"fee"`
	Description  string             `firestore:"description"`
	Children     []BuybackChildItem `firestore:"children"`
}

type BuybackChildItem struct {
	TypeId       int32   `firestore:"type_id"`
	Quantity     float64 `firestore:"quantity"`
	PricePerUnit float64 `firestore:"price_per_unit"`
	Description  string  `firestore:"description"`
}

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
