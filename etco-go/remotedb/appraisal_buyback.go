package remotedb

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
	"github.com/WiggidyW/etco-go/remotedb/implrdb"
)

const (
	B_APPRAISAL_EXPIRES_IN time.Duration = 48 * time.Hour
	B_APPRAISAL_BUF_CAP    int           = 0
)

func init() {
	keys.TypeStrBuybackAppraisal = cache.RegisterType[BuybackAppraisal]("buybackappraisal", B_APPRAISAL_BUF_CAP)
}

type BuybackAppraisal = implrdb.BuybackAppraisal
type BuybackParentItem = implrdb.BuybackParentItem
type BuybackChildItem = implrdb.BuybackChildItem

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
		readBuybackAppraisal,
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
	var cacheLocks []cacheprefetch.ActionOrderedLocks
	if appraisal.CharacterId != nil {
		cacheLocks = []cacheprefetch.ActionOrderedLocks{{
			Locks: []cacheprefetch.ActionLock{
				cacheprefetch.ServerLock(
					keys.CacheKeyUserBuybackAppraisalCodes(
						*appraisal.CharacterId,
					),
					keys.TypeStrUserBuybackAppraisalCodes,
				),
			},
			Child: nil,
		}}
	}
	return appraisalSet(
		x,
		saveBuybackAppraisal,
		keys.TypeStrBuybackAppraisal,
		B_APPRAISAL_EXPIRES_IN,
		appraisal,
		cacheLocks,
	)
}
