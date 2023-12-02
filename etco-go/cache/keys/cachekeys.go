package keys

import (
	"strconv"

	b "github.com/WiggidyW/etco-go-bucket"
	build "github.com/WiggidyW/etco-go/buildconstants"
)

var (
	pfx_nocorp = newPfxKey(
		build.PROGRAM_VERSION,
	)
	pfx_corp = newPfxKey(
		build.PROGRAM_VERSION,
		strconv.Itoa(int(build.CORPORATION_ID)),
		build.DATA_VERSION,
	)
)

// // esi

// jwks
var CacheKeyJWKS = newKey(
	pfx_nocorp,
	"JWKS",
)

// tokens
func CacheKeyAuthToken(refreshToken string) Key {
	return newKey(
		pfx_corp,
		"AuthToken",
		censor(refreshToken),
	)
}
func CacheKeyCorpToken(refreshToken string) Key {
	return newKey(
		pfx_corp,
		"CorpToken",
		censor(refreshToken),
	)
}

func CacheKeyStructureInfoToken(refreshToken string) Key {
	return newKey(
		pfx_corp,
		"StructureInfoToken",
		censor(refreshToken),
	)
}

func CacheKeyMarketsToken(refreshToken string) Key {
	return newKey(
		pfx_corp,
		"MarketsToken",
		censor(refreshToken),
	)
}

// entity info
func CacheKeyAllianceInfo(allianceId int32) Key {
	return newKey(
		pfx_nocorp,
		"AllianceInfo",
		strconv.Itoa(int(allianceId)),
	)
}

func CacheKeyCharacterInfo(characterId int32) Key {
	return newKey(
		pfx_nocorp,
		"CharacterInfo",
		strconv.Itoa(int(characterId)),
	)
}

func CacheKeyCorporationInfo(corporationId int32) Key {
	return newKey(
		pfx_nocorp,
		"CorporationInfo",
		strconv.Itoa(int(corporationId)),
	)
}

// structure info
func CacheKeyStructureInfo(structureId int64) Key {
	return newKey(
		pfx_corp,
		"StructureInfo",
		strconv.Itoa(int(structureId)),
	)
}

// // bucket

// web data
var CacheKeyWebBuybackSystemTypeMapsBuilder = newKey(
	pfx_corp,
	b.OBJNAME_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
)
var CacheKeyWebBuybackBundleKeys = newKey(
	CacheKeyWebBuybackSystemTypeMapsBuilder,
	"BundleKeys",
)
var CacheKeyWebShopLocationTypeMapsBuilder = newKey(
	pfx_corp,
	b.OBJNAME_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
)
var CacheKeyWebShopBundleKeys = newKey(
	CacheKeyWebShopLocationTypeMapsBuilder,
	"BundleKeys",
)
var CacheKeyWebHaulRouteTypeMapsBuilder = newKey(
	pfx_corp,
	b.OBJNAME_WEB_HAUL_ROUTE_TYPE_MAPS_BUILDER,
)
var CacheKeyWebHaulBundleKeys = newKey(
	CacheKeyWebHaulRouteTypeMapsBuilder,
	"BundleKeys",
)
var CacheKeyWebBuybackSystems = newKey(
	pfx_corp,
	b.OBJNAME_WEB_BUYBACK_SYSTEMS,
)
var CacheKeyWebShopLocations = newKey(
	pfx_corp,
	b.OBJNAME_WEB_SHOP_LOCATIONS,
)
var CacheKeyWebHaulRoutes = newKey(
	pfx_corp,
	b.OBJNAME_WEB_HAUL_ROUTES,
)
var CacheKeyWebMarkets = newKey(
	pfx_corp,
	b.OBJNAME_WEB_MARKETS,
)

// build data
var CacheKeyBuildConstData = newKey(
	pfx_corp,
	b.OBJNAME_CONSTANTS_DATA,
)

// auth hash set
func CacheKeyAuthHashSet(domain string) Key {
	return newKey(
		pfx_corp,
		"AuthHashSet",
		domain,
	)
}

// // RemoteDB

// user data
func CacheKeyNSUserData(characterId int32) Key {
	return newKey(
		pfx_corp,
		"UserData",
		strconv.Itoa(int(characterId)),
	)
}
func CacheKeyUserHaulAppraisalCodes(characterId int32) Key {
	return newKey(
		CacheKeyNSUserData(characterId),
		"HaulAppraisalCodes",
	)
}
func CacheKeyUserBuybackAppraisalCodes(characterId int32) Key {
	return newKey(
		CacheKeyNSUserData(characterId),
		"BuybackAppraisalCodes",
	)
}
func CacheKeyUserShopAppraisalCodes(characterId int32) Key {
	return newKey(
		CacheKeyNSUserData(characterId),
		"ShopAppraisalCodes",
	)
}
func CacheKeyUserCancelledPurchase(characterId int32) Key {
	return newKey(
		CacheKeyNSUserData(characterId),
		"CancelledPurchase",
	)
}
func CacheKeyUserMadePurchase(characterId int32) Key {
	return newKey(
		CacheKeyNSUserData(characterId),
		"MadePurchase",
	)
}

// appraisal
func CacheKeyAppraisal(appraisalCode string) Key {
	return newKey(
		pfx_corp,
		"Appraisal",
		appraisalCode,
	)
}

// purchase queue
var CacheKeyRawPurchaseQueue = newKey(
	pfx_corp,
	"RawPurchaseQueue",
)

// previous contracts
var CacheKeyPrevContracts = newKey(
	pfx_corp,
	"PrevContracts",
)

// // Composition
var CacheKeyPurchaseQueue = newKey(
	pfx_corp,
	"PurchaseQueue",
)

func CacheKeyLocationPurchaseQueue(
	locationId int64,
) Key {
	return newKey(
		CacheKeyPurchaseQueue,
		strconv.Itoa(int(locationId)),
	)
}

var CacheKeyNSRawShopAssets = newKey(
	pfx_corp,
	"RawShopAssets",
)

func CacheKeyRawShopAssets(locationId int64) Key {
	return newKey(
		CacheKeyNSRawShopAssets,
		strconv.Itoa(int(locationId)),
	)
}

var CacheKeyNSUnreservedShopAssets = newKey(
	pfx_corp,
	"UnreservedShopAssets",
)

func CacheKeyUnreservedShopAssets(locationId int64) Key {
	return newKey(
		CacheKeyNSUnreservedShopAssets,
		strconv.Itoa(int(locationId)),
	)
}

var CacheKeyNSContracts = newKey(
	pfx_corp,
	"Contracts",
)

var CacheKeyBuybackContracts = newKey(
	CacheKeyNSContracts,
	"Buyback",
)
var CacheKeyShopContracts = newKey(
	CacheKeyNSContracts,
	"Shop",
)
var CacheKeyHaulContracts = newKey(
	CacheKeyNSContracts,
	"Haul",
)

// contractIds are unique
func CacheKeyContractItems(contractId int32) Key {
	return newKey(
		pfx_nocorp,
		"ContractItems",
		strconv.Itoa(int(contractId)),
	)
}

func CacheKeyNSRegionMarketOrders(
	regionId int32,
	typeId int32,
	isBuy bool,
) Key {
	return newKey(
		pfx_nocorp,
		"RegionMarketOrders",
		strconv.Itoa(int(regionId)),
		strconv.Itoa(int(typeId)),
		isBuyStr(isBuy),
	)
}

func CacheKeyRegionMarketOrders(
	nsCacheKey Key,
	locationId int64,
) Key {
	return newKey(
		nsCacheKey,
		strconv.Itoa(int(locationId)),
	)
}

func CacheKeyNSStructureMarketOrders(structureId int64) Key {
	return newKey(
		pfx_nocorp,
		"StructureMarketOrders",
		strconv.Itoa(int(structureId)),
	)
}

func CacheKeyStructureMarketOrders(
	nsCacheKey Key,
	typeId int32,
	isBuy bool,
) Key {
	return newKey(
		nsCacheKey,
		strconv.Itoa(int(typeId)),
		isBuyStr(isBuy),
	)
}

func CacheKeyTokenCharacter(
	app uint8,
	refreshToken string,
) Key {
	return newKey(
		pfx_corp,
		"TokenCharacter",
		strconv.Itoa(int(app)),
		censor(refreshToken),
	)
}
