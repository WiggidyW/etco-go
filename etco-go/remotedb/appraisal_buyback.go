package remotedb

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
	"github.com/WiggidyW/etco-go/proto"
	pr "github.com/WiggidyW/etco-go/protoregistry"
)

const (
	B_APPRAISAL_EXPIRES_IN time.Duration = 48 * time.Hour
	B_APPRAISAL_BUF_CAP    int           = 0
)

func init() {
	keys.TypeStrBuybackAppraisal = cache.RegisterType[BuybackAppraisal]("buybackappraisal", B_APPRAISAL_BUF_CAP)
}

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

func NewBuybackAppraisal(
	rejected bool,
	code string,
	timeStamp time.Time,
	items []BuybackParentItem,
	version string,
	characterId *int32,
	systemId int32,
	price, tax, taxRate, fee, feePerM3 float64,
) BuybackAppraisal {
	return BuybackAppraisal{
		Rejected:    rejected,
		Code:        code,
		Time:        timeStamp,
		Items:       items,
		Version:     version,
		CharacterId: characterId,
		SystemId:    systemId,
		Price:       price,
		Tax:         tax,
		TaxRate:     taxRate,
		Fee:         fee,
		FeePerM3:    feePerM3,
	}
}

func (ba BuybackAppraisal) GetCode() string { return ba.Code }
func (ba BuybackAppraisal) GetCharacterIdVal() (id int32) {
	if ba.CharacterId != nil {
		id = *ba.CharacterId
	}
	return id
}

func (ba BuybackAppraisal) ToProto(
	registry *pr.ProtoRegistry,
) (
	appraisal *proto.BuybackAppraisal,
) {
	return &proto.BuybackAppraisal{
		Rejected:    ba.Rejected,
		Code:        ba.Code,
		Time:        ba.Time.Unix(),
		Items:       proto.P1ToProtoMany(ba.Items, registry),
		Version:     ba.Version,
		CharacterId: ba.GetCharacterIdVal(),
		SystemInfo:  registry.AddSystemById(ba.SystemId),
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
	registry *pr.ProtoRegistry,
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
	registry *pr.ProtoRegistry,
) (
	item *proto.BuybackChildItem,
) {
	return &proto.BuybackChildItem{
		TypeId:              registry.AddTypeById(bci.TypeId),
		QuantityPerParent:   bci.QuantityPerParent,
		PricePerUnit:        bci.PricePerUnit,
		DescriptionStrIndex: registry.Add(bci.Description),
	}
}

func GetBuybackAppraisal(
	x cache.Context,
	code string,
) (
	rep *BuybackAppraisal,
	expires time.Time,
	err error,
) {
	return appraisalGet(
		x,
		client.readBuybackAppraisal,
		keys.TypeStrBuybackAppraisal,
		code,
		B_APPRAISAL_EXPIRES_IN,
	)
}

func SetBuybackAppraisal(
	x cache.Context,
	appraisal BuybackAppraisal,
) (
	err error,
) {
	var cacheLocks []prefetch.CacheActionOrderedLocks
	if appraisal.CharacterId != nil {
		cacheLocks = []prefetch.CacheActionOrderedLocks{
			prefetch.CacheOrderedLocks(
				nil,
				prefetch.ServerCacheLock(
					keys.TypeStrUserBuybackAppraisalCodes,
					keys.CacheKeyUserBuybackAppraisalCodes(
						*appraisal.CharacterId,
					),
				),
			),
		}
	}
	return appraisalSet(
		x,
		client.saveBuybackAppraisal,
		keys.TypeStrBuybackAppraisal,
		B_APPRAISAL_EXPIRES_IN,
		appraisal,
		cacheLocks,
	)
}
