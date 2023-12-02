package appraisal

import (
	"time"

	"github.com/WiggidyW/etco-go/appraisalcode"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/haulsystemids"
	"github.com/WiggidyW/etco-go/items"
	"github.com/WiggidyW/etco-go/market"
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
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
	return create(
		x,
		staticdb.GetHaulRouteInfo,
		func(
			x cache.Context,
			typeId int32,
			quantity int64,
			territoryInfo staticdb.HaulRouteInfo,
		) (
			price remotedb.HaulItem,
			expires time.Time,
			err error,
		) {
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
		TAX_ADD,
		items,
		characterId,
		WrappedHaulSystemIds(haulsystemids.HaulSystemIds{
			Start: startSystemId,
			End:   endSystemId,
		}),
	)
}

func SaveHaulAppraisal(
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
		for _, item := range appraisal.Items {
			if item.PricePerUnit <= 0.0 {
				return appraisal, expires, err
			}
		}
		err = SaveHaulAppraisal(x, appraisal)
	}
	return appraisal, expires, err
}
