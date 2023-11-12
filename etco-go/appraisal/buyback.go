package appraisal

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/fetch/appraisalcode"
	"github.com/WiggidyW/etco-go/market"
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

func CreateBuybackAppraisal[BITEM BasicItem](
	x cache.Context,
	items []BITEM,
	characterId *int32,
	systemId int32,
	includeCode bool,
) (
	appraisal remotedb.BuybackAppraisal,
	expires time.Time,
	err error,
) {
	var codeChar *appraisalcode.CodeChar
	if includeCode {
		codeChar = &appraisalcode.BUYBACK_CHAR
	}
	return create(
		x,
		staticdb.GetBuybackSystemInfo,
		market.GetBuybackPrice,
		remotedb.NewBuybackAppraisal,
		codeChar,
		TAX_SUB,
		items,
		characterId,
		systemId,
	)
}

func SaveBuybackAppraisal(
	x cache.Context,
	appraisal remotedb.BuybackAppraisal,
) (
	err error,
) {
	return remotedb.SetBuybackAppraisal(x, appraisal)
}
