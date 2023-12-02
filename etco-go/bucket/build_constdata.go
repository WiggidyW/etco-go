package bucket

import (
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/proto"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	BUILD_CONST_DATA_BUF_CAP    int           = 0
	BUILD_CONST_DATA_EXPIRES_IN time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrBuildConstData = cache.RegisterType[b.ConstantsData]("buildconstdata", BUILD_CONST_DATA_BUF_CAP)
}

func GetBuildConstData(
	x cache.Context,
) (
	rep b.ConstantsData,
	expires time.Time,
	err error,
) {
	return get(
		x,
		client.ReadConstantsData,
		keys.CacheKeyBuildConstData,
		keys.TypeStrBuildConstData,
		BUILD_CONST_DATA_EXPIRES_IN,
		nil,
	)
}

func ProtoGetBuildConstData(
	x cache.Context,
) (
	rep *proto.CfgConstData,
	expires time.Time,
	err error,
) {
	var buildConstData b.ConstantsData
	buildConstData, expires, err = GetBuildConstData(x)
	if err == nil {
		rep = BuildConstDataToProto(buildConstData)
	}
	return rep, expires, err
}

func SetBuildConstData(
	x cache.Context,
	rep b.ConstantsData,
) (
	err error,
) {
	return set(
		x,
		client.WriteConstantsData,
		keys.CacheKeyBuildConstData,
		keys.TypeStrBuildConstData,
		BUILD_CONST_DATA_EXPIRES_IN,
		rep,
		nil,
	)
}

func ProtoSetBuildConstData(
	x cache.Context,
	rep *proto.CfgConstData,
) (
	err error,
) {
	return SetBuildConstData(x, BuildConstDataFromProto(rep))
}

func BuildConstDataFromProto(
	pbConstData *proto.CfgConstData,
) (
	buildConstData b.ConstantsData,
) {
	if pbConstData == nil {
		return b.ConstantsData{}
	}
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
		DISCORD_CHANNEL:                  &pbConstData.DiscordChannel,
		BUYBACK_CONTRACT_NOTIFICATIONS:   &pbConstData.BuybackContractNotifications,
		SHOP_CONTRACT_NOTIFICATIONS:      &pbConstData.ShopContractNotifications,
		HAUL_CONTRACT_NOTIFICATIONS:      &pbConstData.HaulContractNotifications,
		PURCHASE_NOTIFICATIONS:           &pbConstData.PurchaseNotifications,
	}
}

// TODO: this is despicable
func BuildConstDataToProto(
	buildConstData b.ConstantsData,
) (
	pbConstData *proto.CfgConstData,
) {
	pbConstData = &proto.CfgConstData{}
	if buildConstData.PURCHASE_MAX_ACTIVE != nil {
		pbConstData.PurchaseMaxActive = int32(
			*buildConstData.PURCHASE_MAX_ACTIVE,
		)
	} else {
		pbConstData.PurchaseMaxActive = int32(
			build.PURCHASE_MAX_ACTIVE,
		)
	}
	if buildConstData.MAKE_PURCHASE_COOLDOWN != nil {
		pbConstData.MakePurchaseCooldown = int32(
			buildConstData.MAKE_PURCHASE_COOLDOWN.Nanoseconds() / 1e9,
		)
	} else {
		pbConstData.MakePurchaseCooldown = int32(
			build.MAKE_PURCHASE_COOLDOWN / 1e9,
		)
	}
	if buildConstData.CANCEL_PURCHASE_COOLDOWN != nil {
		pbConstData.CancelPurchaseCooldown = int32(
			buildConstData.MAKE_PURCHASE_COOLDOWN.Nanoseconds() / 1e9,
		)
	} else {
		pbConstData.CancelPurchaseCooldown = int32(
			build.CANCEL_PURCHASE_COOLDOWN / 1e9,
		)
	}
	if buildConstData.CORPORATION_WEB_REFRESH_TOKEN != nil {
		pbConstData.CorporationWebRefreshToken =
			*buildConstData.CORPORATION_WEB_REFRESH_TOKEN
	} else {
		pbConstData.CorporationWebRefreshToken =
			build.CORPORATION_WEB_REFRESH_TOKEN
	}
	if buildConstData.STRUCTURE_INFO_WEB_REFRESH_TOKEN != nil {
		pbConstData.StructureInfoWebRefreshToken =
			*buildConstData.STRUCTURE_INFO_WEB_REFRESH_TOKEN
	} else {
		pbConstData.StructureInfoWebRefreshToken =
			build.STRUCTURE_INFO_WEB_REFRESH_TOKEN
	}
	if buildConstData.DISCORD_CHANNEL != nil {
		pbConstData.DiscordChannel = *buildConstData.DISCORD_CHANNEL
	} else {
		pbConstData.DiscordChannel = build.DISCORD_CHANNEL
	}
	if buildConstData.BUYBACK_CONTRACT_NOTIFICATIONS != nil {
		pbConstData.BuybackContractNotifications =
			*buildConstData.BUYBACK_CONTRACT_NOTIFICATIONS
	} else {
		pbConstData.BuybackContractNotifications =
			build.BUYBACK_CONTRACT_NOTIFICATIONS
	}
	if buildConstData.SHOP_CONTRACT_NOTIFICATIONS != nil {
		pbConstData.ShopContractNotifications =
			*buildConstData.SHOP_CONTRACT_NOTIFICATIONS
	} else {
		pbConstData.ShopContractNotifications =
			build.SHOP_CONTRACT_NOTIFICATIONS
	}
	if buildConstData.HAUL_CONTRACT_NOTIFICATIONS != nil {
		pbConstData.HaulContractNotifications =
			*buildConstData.HAUL_CONTRACT_NOTIFICATIONS
	} else {
		pbConstData.HaulContractNotifications =
			build.HAUL_CONTRACT_NOTIFICATIONS
	}
	if buildConstData.PURCHASE_NOTIFICATIONS != nil {
		pbConstData.PurchaseNotifications =
			*buildConstData.PURCHASE_NOTIFICATIONS
	} else {
		pbConstData.PurchaseNotifications =
			build.PURCHASE_NOTIFICATIONS
	}
	return pbConstData
}
