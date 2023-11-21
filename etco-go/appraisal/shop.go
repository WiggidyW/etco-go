package appraisal

import (
	"time"

	"github.com/WiggidyW/etco-go/appraisalcode"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/items"
	"github.com/WiggidyW/etco-go/market"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
	"github.com/WiggidyW/etco-go/protoregistry"
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

func GetShopAppraisal(
	x cache.Context,
	code string,
) (
	appraisal *remotedb.ShopAppraisal,
	expires time.Time,
	err error,
) {
	return remotedb.GetShopAppraisal(x, code)
}

func GetShopAppraisalCharacterId(
	x cache.Context,
	code string,
) (
	characterId *int32,
	expires time.Time,
	err error,
) {
	var appraisal *remotedb.ShopAppraisal
	appraisal, expires, err = GetShopAppraisal(x, code)
	if err == nil && appraisal != nil {
		characterId = appraisal.CharacterId
	}
	return characterId, expires, err
}

func CreateShopAppraisal[BITEM items.IBasicItem](
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
		status = MPS_None
		err = remotedb.SetShopAppraisal(x, appraisal)
		return appraisalIn, status, err
	} else {
		return userMake(x, appraisal)
	}
}

func CreateSaveShopAppraisal[BITEM items.IBasicItem](
	x cache.Context,
	items []BITEM,
	characterId *int32,
	locationId int64,
) (
	appraisal remotedb.ShopAppraisal,
	status MakePurchaseStatus,
	expires time.Time,
	err error,
) {
	appraisal, expires, err = CreateShopAppraisal(
		x,
		items,
		characterId,
		locationId,
		true,
	)
	if err != nil {
		appraisal, status, err = SaveShopAppraisal(x, appraisal)
	}
	return appraisal, status, expires, err
}

func ProtoGetShopAppraisal(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
	include_items bool,
) (
	appraisal *proto.ShopAppraisal,
	expires time.Time,
	err error,
) {
	var rAppraisal *remotedb.ShopAppraisal
	rAppraisal, expires, err = GetShopAppraisal(x, code)
	if err != nil {
		return nil, expires, protoerr.New(protoerr.SERVER_ERR, err)
	} else if rAppraisal == nil {
		return nil, expires, protoerr.MsgNew(
			protoerr.NOT_FOUND,
			"Shop Appraisal not found",
		)
	} else if !include_items {
		rAppraisal.Items = nil
	}

	var locationInfo *proto.LocationInfo
	locationInfo, expires, err = esi.
		ProtoGetLocationInfoCOV(x, r, rAppraisal.LocationId).
		RecvExpMin(expires)
	if err != nil {
		return nil, expires, err
	}

	return rAppraisal.ToProto(r, locationInfo), expires, nil
}

func ProtoCreateShopAppraisal[BITEM items.IBasicItem](
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	items []BITEM,
	characterId *int32,
	locationId int64,
	save bool,
) (
	appraisal *proto.ShopAppraisal,
	status proto.MakePurchaseStatus,
	expires time.Time,
	err error,
) {
	x, cancel := x.WithCancel()
	defer cancel()
	locationInfoCOV := esi.ProtoGetLocationInfoCOV(x, r, locationId)

	var rAppraisal remotedb.ShopAppraisal
	var rStatus MakePurchaseStatus
	if save {
		rAppraisal, rStatus, expires, err = CreateSaveShopAppraisal(
			x,
			items,
			characterId,
			locationId,
		)
	} else {
		rStatus = MPS_None
		rAppraisal, expires, err = CreateShopAppraisal(
			x,
			items,
			characterId,
			locationId,
			true,
		)
	}
	if err != nil {
		return nil, status, expires, protoerr.New(protoerr.SERVER_ERR, err)
	}

	var locationInfo *proto.LocationInfo
	locationInfo, expires, err = locationInfoCOV.RecvExpMin(expires)
	if err != nil {
		return nil, status, expires, err
	}

	appraisal = rAppraisal.ToProto(r, locationInfo)
	status = rStatus.ToProto()
	return appraisal, status, expires, nil
}
