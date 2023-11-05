package keys

import (
	"fmt"

	b "github.com/WiggidyW/etco-go-bucket"
	build "github.com/WiggidyW/etco-go/buildconstants"
)

var (
	pfx_nocorp = fmt.Sprintf("%s-nocorp", build.DATA_VERSION)
	pfx_corp   = fmt.Sprintf("%s-%d", build.DATA_VERSION, build.CORPORATION_ID)

	CacheKeyJWKS = fmt.Sprintf(
		"%s-JWKS",
		pfx_nocorp,
	)
	CacheKeyContracts = fmt.Sprintf(
		"%s-contracts",
		pfx_corp,
	)
	CacheKeyAllShopAssets = fmt.Sprintf(
		"%s-allshopassets",
		pfx_corp,
	)
	CacheKeyUnreservedShopAssets = fmt.Sprintf(
		"%s-unreservedshopassets",
		pfx_corp,
	)
	CacheKeyPurchaseQueue = fmt.Sprintf(
		"%s-purchasequeue",
		pfx_corp,
	)

	CacheKeyWebBuybackSystemTypeMapsBuilder = fmt.Sprintf(
		"%s-%s",
		pfx_corp,
		b.OBJNAME_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
	)
	CacheKeyWebBuybackBundleKeys = fmt.Sprintf(
		"%s-bundlekeys",
		CacheKeyWebBuybackSystemTypeMapsBuilder,
	)
	CacheKeyWebShopLocationTypeMapsBuilder = fmt.Sprintf(
		"%s-%s",
		pfx_corp,
		b.OBJNAME_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
	)
	CacheKeyWebShopBundleKeys = fmt.Sprintf(
		"%s-bundlekeys",
		CacheKeyWebShopLocationTypeMapsBuilder,
	)
	CacheKeyWebBuybackSystems = fmt.Sprintf(
		"%s-%s",
		pfx_corp,
		b.OBJNAME_WEB_BUYBACK_SYSTEMS,
	)
	CacheKeyWebShopLocations = fmt.Sprintf(
		"%s-%s",
		pfx_corp,
		b.OBJNAME_WEB_SHOP_LOCATIONS,
	)
	CacheKeyWebMarkets = fmt.Sprintf(
		"%s-%s",
		pfx_corp,
		b.OBJNAME_WEB_MARKETS,
	)
	CacheKeyBuildConstData = fmt.Sprintf(
		"%s-%s",
		pfx_corp,
		b.OBJNAME_CONSTANTS_DATA,
	)
)

const (
	IS_BUY_TRUE  string = "b"
	IS_BUY_FALSE string = "s"
)

// contractIds are unique
func CacheKeyContractItems(contractId int32) string {
	return fmt.Sprintf(
		"%s-contractitems-%d",
		pfx_nocorp,
		contractId,
	)
}

func CacheKeyRegionMarket(
	regionId int32,
	typeId int32,
	isBuy bool,
) string {
	return fmt.Sprintf(
		"%s-regionmarket-%d-%d-%s",
		pfx_nocorp,
		regionId,
		typeId,
		isBuyStr(isBuy),
	)
}

func CacheKeyFilterRegionMarket(
	regionId int32,
	typeId int32,
	isBuy bool,
	locationId int64,
) string {
	return fmt.Sprintf(
		"%s-%d",
		CacheKeyRegionMarket(regionId, typeId, isBuy),
		locationId,
	)
}

func CacheKeyStructureMarket(structureId int64) string {
	return fmt.Sprintf(
		"%s-strucmarket-%d",
		pfx_nocorp,
		structureId,
	)
}

func CacheKeyFilterStructureMarket(
	structureId int64,
	typeId int32,
	isBuy bool,
) string {
	return fmt.Sprintf(
		"%s-%d-%s",
		CacheKeyStructureMarket(structureId),
		typeId,
		isBuyStr(isBuy),
	)
}

func CacheKeyCharacterInfo(characterId int32) string {
	return fmt.Sprintf(
		"%s-character-%d",
		pfx_nocorp,
		characterId,
	)
}
func CacheKeyCorporationInfo(corporationId int32) string {
	return fmt.Sprintf(
		"%s-corporation-%d",
		pfx_nocorp,
		corporationId,
	)
}
func CacheKeyAllianceInfo(allianceId int32) string {
	return fmt.Sprintf(
		"%s-alliance-%d",
		pfx_nocorp,
		allianceId,
	)
}

// structure info depends upon docking rights, which are corp-specific
func CacheKeyStructureInfo(structureId int64) string {
	return fmt.Sprintf(
		"%s-strucinfo-%d",
		pfx_corp,
		structureId,
	)
}

func CacheKeyAppraisal(appraisalCode string) string {
	return fmt.Sprintf(
		"%s-appraisal-%s",
		pfx_corp,
		appraisalCode,
	)
}

func CacheKeyUserData(characterId int32) string {
	return fmt.Sprintf(
		"%s-userdata-%d",
		pfx_corp,
		characterId,
	)
}

func CacheKeyUserBuybackAppraisalCodes(characterId int32) string {
	return fmt.Sprintf(
		"%s-b_codes",
		CacheKeyUserData(characterId),
	)
}

func CacheKeyUserShopAppraisalCodes(characterId int32) string {
	return fmt.Sprintf(
		"%s-s_codes",
		CacheKeyUserData(characterId),
	)
}

func CacheKeyUserMadePurchase(characterId int32) string {
	return fmt.Sprintf(
		"%s-m_purchase",
		CacheKeyUserData(characterId),
	)
}

func CacheKeyUserCancelledPurchase(characterId int32) string {
	return fmt.Sprintf(
		"%s-c_purchase",
		CacheKeyUserData(characterId),
	)
}

// NOT SHARED - bucket data is corp-specific
func CacheKeyAuthHashSet(domain string) string {
	return fmt.Sprintf(
		"%s-authhashset-%s",
		pfx_corp,
		domain,
	)
}

// // util
func isBuyStr(isBuy bool) string {
	if isBuy {
		return IS_BUY_TRUE
	} else {
		return IS_BUY_FALSE
	}
}
