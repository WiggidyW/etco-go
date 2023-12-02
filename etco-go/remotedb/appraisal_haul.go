package remotedb

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
	"github.com/WiggidyW/etco-go/haulsystemids"
	"github.com/WiggidyW/etco-go/remotedb/implrdb"
)

const (
	H_APPRAISAL_EXPIRES_IN time.Duration = 48 * time.Hour
	H_APPRAISAL_BUF_CAP    int           = 0
)

func init() {
	keys.TypeStrHaulAppraisal = cache.RegisterType[HaulAppraisal]("haulappraisal", H_APPRAISAL_BUF_CAP)
}

type HaulAppraisal = implrdb.HaulAppraisal
type HaulItem = implrdb.HaulItem

func NewHaulAppraisal(
	rejected bool,
	code string,
	timeStamp time.Time,
	items []HaulItem,
	version string,
	characterId *int32,
	systemIds haulsystemids.HaulSystemIds,
	price, tax, taxRate, fee, feePerM3 float64,
) HaulAppraisal {
	return HaulAppraisal{
		Rejected:      rejected,
		Code:          code,
		Time:          timeStamp,
		Items:         items,
		Version:       version,
		CharacterId:   characterId,
		StartSystemId: systemIds.Start,
		EndSystemId:   systemIds.End,
		Price:         price,
		Tax:           tax,
		TaxRate:       taxRate,
		Fee:           fee,
		FeePerM3:      feePerM3,
	}
}

func GetHaulAppraisal(
	x cache.Context,
	code string,
) (
	rep *HaulAppraisal,
	expires time.Time,
	err error,
) {
	return appraisalGet(
		x,
		readHaulAppraisal,
		keys.TypeStrHaulAppraisal,
		code,
		H_APPRAISAL_EXPIRES_IN,
	)
}

func SetHaulAppraisal(
	x cache.Context,
	appraisal HaulAppraisal,
) (
	err error,
) {
	var cacheLocks []cacheprefetch.ActionOrderedLocks
	if appraisal.CharacterId != nil {
		cacheLocks = []cacheprefetch.ActionOrderedLocks{{
			Locks: []cacheprefetch.ActionLock{
				cacheprefetch.ServerLock(
					keys.CacheKeyUserHaulAppraisalCodes(
						*appraisal.CharacterId,
					),
					keys.TypeStrUserHaulAppraisalCodes,
				),
			},
		}}
	}
	return appraisalSet(
		x,
		saveHaulAppraisal,
		keys.TypeStrHaulAppraisal,
		H_APPRAISAL_EXPIRES_IN,
		appraisal,
		cacheLocks,
	)
}
