package appraisal

import (
	"time"

	"github.com/WiggidyW/etco-go/appraisalcode"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/haulsystemids"
	"github.com/WiggidyW/etco-go/items"
	"github.com/WiggidyW/etco-go/market"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
	"github.com/WiggidyW/etco-go/protoregistry"
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	HAUL_ITEM_FALLBACK_PRICE_EXPIRES time.Duration = 1 * time.Hour
)

func GetHaulAppraisal(
	x cache.Context,
	code string,
) (
	appraisal *remotedb.HaulAppraisal,
	expires time.Time,
	err error,
) {
	return remotedb.GetHaulAppraisal(x, code)
}

func ProtoGetHaulAppraisal(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
	include_items bool,
) (
	appraisal *proto.HaulAppraisal,
	expires time.Time,
	err error,
) {
	var rAppraisal *remotedb.HaulAppraisal
	rAppraisal, _, err = GetHaulAppraisal(x, code)
	if err != nil {
		return nil, expires, err
	} else if rAppraisal == nil {
		return nil, expires, protoerr.MsgNew(
			protoerr.NOT_FOUND,
			"Haul Appraisal not found",
		)
	} else if !include_items {
		rAppraisalCopy := *rAppraisal
		rAppraisal = &rAppraisalCopy
		rAppraisal.Items = nil
	}

	return rAppraisal.ToProto(r), expires, nil
}

func GetHaulAppraisalCharacterId(
	x cache.Context,
	code string,
) (
	characterId *int32,
	expires time.Time,
	err error,
) {
	var appraisal *remotedb.HaulAppraisal
	appraisal, expires, err = GetHaulAppraisal(x, code)
	if err == nil && appraisal != nil {
		characterId = appraisal.CharacterId
	}
	return characterId, expires, err
}

type WrappedHaulSystemIds haulsystemids.HaulSystemIds

// only used for hashing
func (w WrappedHaulSystemIds) ToInt64() int64 {
	return int64(w.Start)<<32 | int64(w.End)
}
func (w WrappedHaulSystemIds) Unwrap() haulsystemids.HaulSystemIds {
	return haulsystemids.HaulSystemIds(w)
}

func CreateHaulAppraisal[BITEM items.IBasicItem](
	x cache.Context,
	items []BITEM,
	characterId *int32,
	startSystemId, endSystemId int32,
	includeCode bool,
	fallbackPricePerUnit map[int32]float64,
) (
	appraisal remotedb.HaulAppraisal,
	expires time.Time,
	err error,
) {
	var codeChar *appraisalcode.CodeChar
	if includeCode {
		codeChar = &appraisalcode.HAUL_CHAR
	}
	var routeInfo staticdb.HaulRouteInfo
	var m3Fee float64
	routeInfo, m3Fee, appraisal, expires, err = create(
		x,
		staticdb.GetHaulRouteInfo,
		func(
			x cache.Context,
			typeId int32,
			quantity int64,
			territoryInfo staticdb.HaulRouteInfo,
		) (remotedb.HaulItem, time.Time, error) {
			return market.GetHaulPrice(
				x,
				typeId,
				quantity,
				territoryInfo,
				fallbackPricePerUnit,
			)
		},
		remotedb.NewHaulAppraisal,
		codeChar,
		PTAX_NONE,
		items,
		characterId,
		WrappedHaulSystemIds(haulsystemids.HaulSystemIds{
			Start: startSystemId,
			End:   endSystemId,
		}),
	)

	// set rejected if any item is rejected
	if !appraisal.Rejected {
		for _, item := range appraisal.Items {
			if item.PricePerUnit <= 0.0 {
				appraisal.Rejected = true
				break
			}
		}
	}

	// add reward and tax
	appraisal = HaulAppraisalWithRewardAndTax(appraisal, routeInfo, m3Fee)

	return appraisal, expires, err
}

func ProtoCreateHaulAppraisal[BITEM items.IBasicItem](
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	items []BITEM,
	characterId *int32,
	startSystemId, endSystemId int32,
	save bool,
	fallbackPricePerUnit map[int32]float64,
) (
	appraisal *proto.HaulAppraisal,
	expires time.Time,
	err error,
) {
	var rAppraisal remotedb.HaulAppraisal
	if save {
		rAppraisal, expires, err = CreateSaveHaulAppraisal(
			x,
			items,
			characterId,
			startSystemId, endSystemId,
			fallbackPricePerUnit,
		)
	} else {
		rAppraisal, expires, err = CreateHaulAppraisal(
			x,
			items,
			characterId,
			startSystemId, endSystemId,
			false,
			fallbackPricePerUnit,
		)
	}
	if err != nil {
		return nil, expires, err
	} else {
		return rAppraisal.ToProto(r), expires, nil
	}
}

func saveHaulAppraisal(
	x cache.Context,
	appraisal remotedb.HaulAppraisal,
) (
	err error,
) {
	return remotedb.SetHaulAppraisal(x, appraisal)
}

// Only saves if no items are rejected
func CreateSaveHaulAppraisal[BITEM items.IBasicItem](
	x cache.Context,
	items []BITEM,
	characterId *int32,
	startSystemId, endSystemId int32,
	fallbackPricePerUnit map[int32]float64,
) (
	appraisal remotedb.HaulAppraisal,
	expires time.Time,
	err error,
) {
	appraisal, expires, err = CreateHaulAppraisal(
		x,
		items,
		characterId,
		startSystemId, endSystemId,
		true,
		fallbackPricePerUnit,
	)
	if err == nil && !appraisal.Rejected && appraisal.Price > 0.0 {
		err = saveHaulAppraisal(x, appraisal)
	}
	return appraisal, expires, err
}

func HaulAppraisalWithRewardAndTax(
	appraisalIn remotedb.HaulAppraisal,
	routeInfo staticdb.HaulRouteInfo,
	m3Reward float64,
) (
	appraisal remotedb.HaulAppraisal,
) {
	appraisal = appraisalIn
	appraisal.CollateralRate = routeInfo.CollateralRate

	// compute reward, ignoring min and max reward
	collateralReward := appraisal.Price * appraisal.CollateralRate
	rewardKind := remotedb.HRKInvalid
	if routeInfo.RewardStrategy == b.HRRSSum {
		rewardKind = remotedb.HRKSum
		appraisal.Reward = m3Reward + collateralReward
	} else if (routeInfo.RewardStrategy == b.HRRSGreaterOf &&
		m3Reward > collateralReward) ||
		(routeInfo.RewardStrategy == b.HRRSLesserOf &&
			m3Reward < collateralReward) {
		rewardKind = remotedb.HRKM3Fee
		appraisal.Reward = m3Reward
	} else /* if (routeInfo.RewardStrategy == b.HRRSGreaterOf &&
	   collateralReward >= m3Reward) ||
	   (routeInfo.RewardStrategy == b.HRRSLesserOf &&
	   collateralReward <= m3Reward) */{
		rewardKind = remotedb.HRKCollateral
		appraisal.Reward = collateralReward
	}

	// apply min and max reward
	if appraisal.Reward < routeInfo.MinReward {
		switch rewardKind {
		case remotedb.HRKCollateral:
			rewardKind = remotedb.HRKMinRewardCollateral
		case remotedb.HRKM3Fee:
			rewardKind = remotedb.HRKMinRewardM3Fee
		case remotedb.HRKSum:
			rewardKind = remotedb.HRKMinRewardSum
		}
		appraisal.Reward = routeInfo.MinReward
	} else if appraisal.Reward > routeInfo.MaxReward {
		switch rewardKind {
		case remotedb.HRKCollateral:
			rewardKind = remotedb.HRKMaxRewardCollateral
		case remotedb.HRKM3Fee:
			rewardKind = remotedb.HRKMaxRewardM3Fee
		case remotedb.HRKSum:
			rewardKind = remotedb.HRKMaxRewardSum
		}
		appraisal.Reward = routeInfo.MaxReward
	}

	// set tax and add tax to reward
	appraisal.Tax = appraisal.Reward * routeInfo.TaxRate
	appraisal.Reward += appraisal.Tax

	// set reward kind (cast it to uint8)
	appraisal.RewardKind = rewardKind.Uint8()

	return appraisal
}
