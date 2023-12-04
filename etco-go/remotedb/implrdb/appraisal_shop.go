package implrdb

import (
	"time"

	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"
)

type ShopAppraisal struct {
	Rejected bool `firestore:"rejected,omitempty"`

	// ignored during reading (used as doc id instead of field)
	// technically, if you're reading, you must already know it
	Code string `firestore:"-"`

	Time        time.Time  `firestore:"time"`
	Items       []ShopItem `firestore:"items"`
	Version     string     `firestore:"version"`
	CharacterId *int32     `firestore:"character_id"`
	LocationId  int64      `firestore:"location_id"`
	Price       float64    `firestore:"price"`
	Tax         float64    `firestore:"tax,omitempty"`
	TaxRate     float64    `firestore:"tax_rate,omitempty"`
}

func (sa ShopAppraisal) GetCode() string { return sa.Code }
func (sa ShopAppraisal) GetCharacterIdVal() (id int32) {
	if sa.CharacterId != nil {
		id = *sa.CharacterId
	}
	return id
}

func (sa ShopAppraisal) ToProto(
	r *protoregistry.ProtoRegistry,
	locationInfo *proto.LocationInfo,
) (
	appraisal *proto.ShopAppraisal,
) {
	return &proto.ShopAppraisal{
		Rejected:     sa.Rejected,
		Code:         sa.Code,
		Time:         sa.Time.Unix(),
		Items:        proto.P1ToProtoMany(sa.Items, r),
		Version:      sa.Version,
		CharacterId:  sa.GetCharacterIdVal(),
		LocationInfo: locationInfo,
		Price:        sa.Price,
		Tax:          sa.Tax,
		TaxRate:      sa.TaxRate,
	}
}

type ShopItem struct {
	TypeId       int32   `firestore:"type_id"`
	Quantity     int64   `firestore:"quantity"`
	PricePerUnit float64 `firestore:"price_per_unit"`
	Description  string  `firestore:"description"`
}

func (si ShopItem) GetTypeId() int32         { return si.TypeId }
func (si ShopItem) GetQuantity() int64       { return si.Quantity }
func (si ShopItem) GetPricePerUnit() float64 { return si.PricePerUnit }
func (si ShopItem) GetDescription() string   { return si.Description }
func (si ShopItem) GetFeePerUnit() float64   { return 0.0 }
func (si ShopItem) GetChildrenLength() int   { return 0 }

func (si ShopItem) ToProto(
	r *protoregistry.ProtoRegistry,
) (
	item *proto.ShopItem,
) {
	return &proto.ShopItem{
		TypeId:              r.AddTypeById(si.TypeId),
		Quantity:            si.Quantity,
		PricePerUnit:        si.PricePerUnit,
		DescriptionStrIndex: r.Add(si.Description),
	}
}
