package appraisal

import (
	"time"

	"github.com/WiggidyW/etco-go/appraisalcode"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/market"
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

func CreateShopAppraisal[BITEM BasicItem](
	x cache.Context,
	items []BITEM,
	characterId *int32,
	locationId int64,
	includeCode bool,
) (
	appraisal remotedb.ShopAppraisal,
	expires time.Time,
	err error,
) {
	var codeChar *appraisalcode.CodeChar
	if includeCode {
		codeChar = &appraisalcode.BUYBACK_CHAR
	}
	return create(
		x,
		staticdb.GetShopLocationInfo,
		market.GetShopPrice,
		remotedb.NewShopAppraisal,
		codeChar,
		TAX_ADD,
		items,
		characterId,
		locationId,
	)
}

func SaveShopAppraisal(
	x cache.Context,
	appraisalIn remotedb.ShopAppraisal,
) (
	appraisal remotedb.ShopAppraisal,
	status MakePurchaseStatus,
	err error,
) {
	if appraisal.CharacterId == nil {
		status = MPS_Success
		err = remotedb.SetShopAppraisal(x, appraisal)
		return appraisalIn, status, err
	} else {
		return userMake(x, appraisal)
	}
}
