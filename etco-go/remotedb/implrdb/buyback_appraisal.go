package implrdb

import (
	"time"

	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"
)

type BuybackAppraisal struct {
	Rejected bool `firestore:"rejected,omitempty"`

	// ignored during reading (used as doc id instead of field)
	// technically, if you're reading, you must already know it
	Code string `firestore:"-"`

	Time        time.Time           `firestore:"time"`
	Items       []BuybackParentItem `firestore:"items"`
	Version     string              `firestore:"version"`
	CharacterId *int32              `firestore:"character_id"`
	SystemId    int32               `firestore:"system_id"`
	Price       float64             `firestore:"price"`
	Tax         float64             `firestore:"tax,omitempty"`
	TaxRate     float64             `firestore:"tax_rate,omitempty"`
	Fee         float64             `firestore:"fee,omitempty"`
	FeePerM3    float64             `firestore:"fee_per_m3,omitempty"`
}

func (ba BuybackAppraisal) GetCode() string { return ba.Code }
func (ba BuybackAppraisal) GetCharacterIdVal() (id int32) {
	if ba.CharacterId != nil {
		id = *ba.CharacterId
	}
	return id
}

func (ba BuybackAppraisal) ToProto(
	r *protoregistry.ProtoRegistry,
) (
	appraisal *proto.BuybackAppraisal,
) {
	return &proto.BuybackAppraisal{
		Rejected:    ba.Rejected,
		Code:        ba.Code,
		Time:        ba.Time.Unix(),
		Items:       proto.P1ToProtoMany(ba.Items, r),
		Version:     ba.Version,
		CharacterId: ba.GetCharacterIdVal(),
		SystemInfo:  r.AddSystemById(ba.SystemId),
		Price:       ba.Price,
		Tax:         ba.Tax,
		TaxRate:     ba.TaxRate,
		Fee:         ba.Fee,
		FeePerM3:    ba.FeePerM3,
	}
}

type BuybackParentItem struct {
	TypeId       int32              `firestore:"type_id"`
	Quantity     int64              `firestore:"quantity"`
	PricePerUnit float64            `firestore:"price_per_unit"`
	Description  string             `firestore:"description"`
	FeePerUnit   float64            `firestore:"fee,omitempty"`
	Children     []BuybackChildItem `firestore:"children"`
}

func (bpi BuybackParentItem) GetTypeId() int32         { return bpi.TypeId }
func (bpi BuybackParentItem) GetQuantity() int64       { return bpi.Quantity }
func (bpi BuybackParentItem) GetPricePerUnit() float64 { return bpi.PricePerUnit }
func (bpi BuybackParentItem) GetDescription() string   { return bpi.Description }
func (bpi BuybackParentItem) GetFeePerUnit() float64   { return bpi.FeePerUnit }
func (bpi BuybackParentItem) GetChildrenLength() int   { return len(bpi.Children) }

func (bpi BuybackParentItem) ToProto(
	registry *protoregistry.ProtoRegistry,
) (
	item *proto.BuybackParentItem,
) {
	return &proto.BuybackParentItem{
		TypeId:              registry.AddTypeById(bpi.TypeId),
		Quantity:            bpi.Quantity,
		PricePerUnit:        bpi.PricePerUnit,
		DescriptionStrIndex: registry.Add(bpi.Description),
		FeePerUnit:          bpi.FeePerUnit,
		Children:            proto.P1ToProtoMany(bpi.Children, registry),
	}
}

type BuybackChildItem struct {
	TypeId            int32   `firestore:"type_id"`
	QuantityPerParent float64 `firestore:"quantity_per_parent"`
	PricePerUnit      float64 `firestore:"price_per_unit"`
	Description       string  `firestore:"description"`
}

func (bci BuybackChildItem) ToProto(
	r *protoregistry.ProtoRegistry,
) (
	item *proto.BuybackChildItem,
) {
	return &proto.BuybackChildItem{
		TypeId:              r.AddTypeById(bci.TypeId),
		QuantityPerParent:   bci.QuantityPerParent,
		PricePerUnit:        bci.PricePerUnit,
		DescriptionStrIndex: r.Add(bci.Description),
	}
}
