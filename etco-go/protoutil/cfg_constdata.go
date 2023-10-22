package protoutil

import (
	"time"

	b "github.com/WiggidyW/etco-go-bucket"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/proto"
)

// This demonstrates for me just how verbose this language is.
func NewPBConstData(
	rConstData b.ConstantsData,
) (
	pbConstData *proto.ConstData,
) {
	pbConstData = &proto.ConstData{}
	if rConstData.PURCHASE_MAX_ACTIVE != nil {
		pbConstData.PurchaseMaxActive = int32(
			*rConstData.PURCHASE_MAX_ACTIVE,
		)
	} else {
		pbConstData.PurchaseMaxActive = int32(
			build.PURCHASE_MAX_ACTIVE,
		)
	}
	if rConstData.MAKE_PURCHASE_COOLDOWN != nil {
		pbConstData.MakePurchaseCooldown = int32(
			rConstData.MAKE_PURCHASE_COOLDOWN.Nanoseconds() / 1e9,
		)
	} else {
		pbConstData.MakePurchaseCooldown = int32(
			build.MAKE_PURCHASE_COOLDOWN / 1e9,
		)
	}
	if rConstData.CANCEL_PURCHASE_COOLDOWN != nil {
		pbConstData.CancelPurchaseCooldown = int32(
			rConstData.MAKE_PURCHASE_COOLDOWN.Nanoseconds() / 1e9,
		)
	} else {
		pbConstData.CancelPurchaseCooldown = int32(
			build.CANCEL_PURCHASE_COOLDOWN / 1e9,
		)
	}
	if rConstData.CORPORATION_WEB_REFRESH_TOKEN != nil {
		pbConstData.CorporationWebRefreshToken =
			*rConstData.CORPORATION_WEB_REFRESH_TOKEN
	} else {
		pbConstData.CorporationWebRefreshToken =
			build.CORPORATION_WEB_REFRESH_TOKEN
	}
	if rConstData.STRUCTURE_INFO_WEB_REFRESH_TOKEN != nil {
		pbConstData.StructureInfoWebRefreshToken =
			*rConstData.STRUCTURE_INFO_WEB_REFRESH_TOKEN
	} else {
		pbConstData.StructureInfoWebRefreshToken =
			build.STRUCTURE_INFO_WEB_REFRESH_TOKEN
	}
	return pbConstData
}

func NewRConstData(
	pbConstData *proto.ConstData,
) (
	rConstData b.ConstantsData,
) {
	PURCHASE_MAX_ACTIVE := int(pbConstData.PurchaseMaxActive)
	MAKE_PURCHASE_COOLDOWN := time.Duration(
		int64(pbConstData.MakePurchaseCooldown) * 1e9,
	)
	CANCEL_PURCHASE_COOLDOWN := time.Duration(
		int64(pbConstData.CancelPurchaseCooldown) * 1e9,
	)
	return b.ConstantsData{
		PURCHASE_MAX_ACTIVE:              &PURCHASE_MAX_ACTIVE,
		MAKE_PURCHASE_COOLDOWN:           &MAKE_PURCHASE_COOLDOWN,
		CANCEL_PURCHASE_COOLDOWN:         &CANCEL_PURCHASE_COOLDOWN,
		CORPORATION_WEB_REFRESH_TOKEN:    &pbConstData.CorporationWebRefreshToken,
		STRUCTURE_INFO_WEB_REFRESH_TOKEN: &pbConstData.StructureInfoWebRefreshToken,
	}
}
