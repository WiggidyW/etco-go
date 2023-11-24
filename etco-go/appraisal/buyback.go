package appraisal

import (
	"time"

	"github.com/WiggidyW/etco-go/appraisalcode"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/items"
	"github.com/WiggidyW/etco-go/market"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
	"github.com/WiggidyW/etco-go/protoregistry"
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

func GetBuybackAppraisal(
	x cache.Context,
	code string,
) (
	appraisal *remotedb.BuybackAppraisal,
	expires time.Time,
	err error,
) {
	return remotedb.GetBuybackAppraisal(x, code)
}

func ProtoGetBuybackAppraisal(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
	include_items bool,
) (
	appraisal *proto.BuybackAppraisal,
	expires time.Time,
	err error,
) {
	var rAppraisal *remotedb.BuybackAppraisal
	rAppraisal, _, err = GetBuybackAppraisal(x, code)
	if err != nil {
		return nil, expires, err
	} else if rAppraisal == nil {
		return nil, expires, protoerr.MsgNew(
			protoerr.NOT_FOUND,
			"Buyback Appraisal not found",
		)
	} else if !include_items {
		rAppraisal.Items = nil
	}

	return rAppraisal.ToProto(r), expires, nil
}

func GetBuybackAppraisalCharacterId(
	x cache.Context,
	code string,
) (
	characterId *int32,
	expires time.Time,
	err error,
) {
	var appraisal *remotedb.BuybackAppraisal
	appraisal, expires, err = GetBuybackAppraisal(x, code)
	if err == nil && appraisal != nil {
		characterId = appraisal.CharacterId
	}
	return characterId, expires, err
}

func CreateBuybackAppraisal[BITEM items.IBasicItem](
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

func ProtoCreateBuybackAppraisal[BITEM items.IBasicItem](
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	items []BITEM,
	characterId *int32,
	systemId int32,
	save bool,
) (
	appraisal *proto.BuybackAppraisal,
	expires time.Time,
	err error,
) {
	var rAppraisal remotedb.BuybackAppraisal
	if save {
		rAppraisal, expires, err = CreateSaveBuybackAppraisal(
			x,
			items,
			characterId,
			systemId,
		)
	} else {
		rAppraisal, expires, err = CreateBuybackAppraisal(
			x,
			items,
			characterId,
			systemId,
			false,
		)
	}
	if err != nil {
		return nil, expires, err
	} else {
		return rAppraisal.ToProto(r), expires, nil
	}
}

func SaveBuybackAppraisal(
	x cache.Context,
	appraisal remotedb.BuybackAppraisal,
) (
	err error,
) {
	return remotedb.SetBuybackAppraisal(x, appraisal)
}

func CreateSaveBuybackAppraisal[BITEM items.IBasicItem](
	x cache.Context,
	items []BITEM,
	characterId *int32,
	systemId int32,
) (
	appraisal remotedb.BuybackAppraisal,
	expires time.Time,
	err error,
) {
	appraisal, expires, err = CreateBuybackAppraisal(
		x,
		items,
		characterId,
		systemId,
		true,
	)
	if err == nil && !appraisal.Rejected && appraisal.Price > 0.0 {
		err = SaveBuybackAppraisal(x, appraisal)
	}
	return appraisal, expires, err
}
