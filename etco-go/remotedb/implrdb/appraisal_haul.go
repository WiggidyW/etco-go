package implrdb

import (
	"time"
)

type HaulAppraisal struct {
	Rejected bool `firestore:"rejected,omitempty"`

	// ignored during reading (used as doc id instead of field)
	// technically, if you're reading, you must already know it
	Code string `firestore:"-"`

	Time          time.Time  `firestore:"time"`
	Items         []HaulItem `firestore:"items"`
	Version       string     `firestore:"version"`
	CharacterId   *int32     `firestore:"character_id"`
	StartSystemId int32      `firestore:"start_system_id"`
	EndSystemId   int32      `firestore:"end_system_id"`
	Price         float64    `firestore:"price"`
	Tax           float64    `firestore:"tax,omitempty"`
	TaxRate       float64    `firestore:"tax_rate,omitempty"`
	Fee           float64    `firestore:"fee,omitempty"`
	FeePerM3      float64    `firestore:"fee_per_m3,omitempty"`
}

func (ha HaulAppraisal) GetCode() string { return ha.Code }
func (ha HaulAppraisal) GetCharacterIdVal() (id int32) {
	if ha.CharacterId != nil {
		id = *ha.CharacterId
	}
	return id
}

type HaulItem struct {
	TypeId       int32   `firestore:"type_id"`
	Quantity     int64   `firestore:"quantity"`
	PricePerUnit float64 `firestore:"price_per_unit"`
	Description  string  `firestore:"description"`
	FeePerUnit   float64 `firestore:"fee,omitempty"`
}

func (hi HaulItem) GetTypeId() int32         { return hi.TypeId }
func (hi HaulItem) GetQuantity() int64       { return hi.Quantity }
func (hi HaulItem) GetPricePerUnit() float64 { return hi.PricePerUnit }
func (hi HaulItem) GetDescription() string   { return hi.Description }
func (hi HaulItem) GetFeePerUnit() float64   { return hi.FeePerUnit }
func (hi HaulItem) GetChildrenLength() int   { return 0 }
