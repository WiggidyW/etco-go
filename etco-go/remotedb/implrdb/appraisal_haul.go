package implrdb

import (
	"time"

	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"
)

type HaulAppraisalRewardKind uint8

const (
	HRKInvalid             HaulAppraisalRewardKind = 0
	HRKCollateral          HaulAppraisalRewardKind = 101 // reward == collateral * collateralRate
	HRKM3Fee               HaulAppraisalRewardKind = 102 // reward == m3Fee * volume
	HRKSum                 HaulAppraisalRewardKind = 103 // reward == (m3Fee * volume) + (collateral * collateralRate)
	HRKMinRewardCollateral HaulAppraisalRewardKind = 104 // reward would've been collateral, but min reward >
	HRKMinRewardM3Fee      HaulAppraisalRewardKind = 105 // reward would've been m3Fee, but min reward >
	HRKMinRewardSum        HaulAppraisalRewardKind = 106 // reward would've been sum, but min reward >
	HRKMaxRewardCollateral HaulAppraisalRewardKind = 107 // reward would've been collateral, but max reward <
	HRKMaxRewardM3Fee      HaulAppraisalRewardKind = 108 // reward would've been m3Fee, but max reward <
	HRKMaxRewardSum        HaulAppraisalRewardKind = 109 // reward would've been sum, but max reward <
)

func (rk HaulAppraisalRewardKind) Uint8() uint8 { return uint8(rk) }

func (rk HaulAppraisalRewardKind) ToProto() proto.HaulRewardKind {
	switch rk {
	case HRKCollateral:
		return proto.HaulRewardKind_HRK_COLLATERAL
	case HRKM3Fee:
		return proto.HaulRewardKind_HRK_M3FEE
	case HRKSum:
		return proto.HaulRewardKind_HRK_SUM
	case HRKMinRewardCollateral:
		return proto.HaulRewardKind_HRK_MIN_REWARD_COLLATERAL
	case HRKMinRewardM3Fee:
		return proto.HaulRewardKind_HRK_MIN_REWARD_M3FEE
	case HRKMinRewardSum:
		return proto.HaulRewardKind_HRK_MIN_REWARD_SUM
	case HRKMaxRewardCollateral:
		return proto.HaulRewardKind_HRK_MAX_REWARD_COLLATERAL
	case HRKMaxRewardM3Fee:
		return proto.HaulRewardKind_HRK_MAX_REWARD_M3FEE
	case HRKMaxRewardSum:
		return proto.HaulRewardKind_HRK_MAX_REWARD_SUM
	default:
		return proto.HaulRewardKind_HRK_INVALID
	}
}

type HaulAppraisal struct {
	Rejected bool `firestore:"rejected,omitempty"`

	// ignored during reading (used as doc id instead of field)
	// technically, if you're reading, you must already know it
	Code string `firestore:"-"`

	Time           time.Time  `firestore:"time"`
	Items          []HaulItem `firestore:"items"`
	Version        string     `firestore:"version"`
	CharacterId    *int32     `firestore:"character_id"`
	StartSystemId  int32      `firestore:"start_system_id"`
	EndSystemId    int32      `firestore:"end_system_id"`
	Price          float64    `firestore:"price"`
	Tax            float64    `firestore:"tax,omitempty"`
	TaxRate        float64    `firestore:"tax_rate,omitempty"`
	FeePerM3       float64    `firestore:"fee_per_m3,omitempty"`
	CollateralRate float64    `firestore:"collateral_rate,omitempty"`
	Reward         float64    `firestore:"reward"`
	RewardKind     uint8      `firestore:"reward_kind"`
}

func (ha HaulAppraisal) GetCode() string { return ha.Code }
func (ha HaulAppraisal) GetCharacterIdVal() (id int32) {
	if ha.CharacterId != nil {
		id = *ha.CharacterId
	}
	return id
}
func (ha HaulAppraisal) GetRewardKind() HaulAppraisalRewardKind {
	switch ha.RewardKind {
	case 101:
		return HRKCollateral
	case 102:
		return HRKM3Fee
	case 103:
		return HRKSum
	case 104:
		return HRKMinRewardCollateral
	case 105:
		return HRKMinRewardM3Fee
	case 106:
		return HRKMinRewardSum
	case 107:
		return HRKMaxRewardCollateral
	case 108:
		return HRKMaxRewardM3Fee
	case 109:
		return HRKMaxRewardSum
	default:
		return HRKInvalid
	}
}

func (ha HaulAppraisal) ToProto(
	r *protoregistry.ProtoRegistry,
) (
	appraisal *proto.HaulAppraisal,
) {
	return &proto.HaulAppraisal{
		Rejected:    ha.Rejected,
		Code:        ha.Code,
		Time:        ha.Time.Unix(),
		Items:       proto.P1ToProtoMany(ha.Items, r),
		Version:     ha.Version,
		CharacterId: ha.GetCharacterIdVal(),
		RouteInfo: &proto.HaulRouteInfo{
			StartSystemInfo: r.AddSystemById(ha.StartSystemId),
			EndSystemInfo:   r.AddSystemById(ha.EndSystemId),
		},
		Price:          ha.Price,
		Tax:            ha.Tax,
		TaxRate:        ha.TaxRate,
		FeePerM3:       ha.FeePerM3,
		CollateralRate: ha.CollateralRate,
		Reward:         ha.Reward,
		RewardKind:     ha.GetRewardKind().ToProto(),
	}
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

func (hi HaulItem) ToProto(
	r *protoregistry.ProtoRegistry,
) (
	item *proto.HaulItem,
) {
	return &proto.HaulItem{
		TypeId:              r.AddTypeById(hi.TypeId),
		Quantity:            hi.Quantity,
		PricePerUnit:        hi.PricePerUnit,
		DescriptionStrIndex: r.Add(hi.Description),
		FeePerUnit:          hi.FeePerUnit,
	}
}
