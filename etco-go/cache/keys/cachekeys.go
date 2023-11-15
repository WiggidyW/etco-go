package keys

import (
	"fmt"
	"strconv"

	b "github.com/WiggidyW/etco-go-bucket"
	build "github.com/WiggidyW/etco-go/buildconstants"
)

var (
	pfx_nocorp = fmt.Sprintf("%s-nocorp", build.DATA_VERSION)
	pfx_corp   = fmt.Sprintf("%s-%d", build.DATA_VERSION, build.CORPORATION_ID)
)

func newCacheKey(parent string, keyParts ...string) string {
	key := parent
	for _, keyPart := range keyParts {
		key += "-"
		key += keyPart
	}
	return key
}

// // esi

// jwks
var CacheKeyJWKS = newCacheKey(
	pfx_nocorp,
	"JWKS",
)

// tokens
func CacheKeyAuthToken(refreshToken string) string {
	return newCacheKey(
		pfx_corp,
		"AuthToken",
		refreshToken,
	)
}
func CacheKeyCorpToken(refreshToken string) string {
	return newCacheKey(
		pfx_corp,
		"CorpToken",
		refreshToken,
	)
}

func CacheKeyStructureInfoToken(refreshToken string) string {
	return newCacheKey(
		pfx_corp,
		"StructureInfoToken",
		refreshToken,
	)
}

func CacheKeyMarketsToken(refreshToken string) string {
	return newCacheKey(
		pfx_corp,
		"MarketsToken",
		refreshToken,
	)
}

// entity info
func CacheKeyAllianceInfo(allianceId int32) string {
	return newCacheKey(
		pfx_nocorp,
		"AllianceInfo",
		strconv.Itoa(int(allianceId)),
	)
}

func CacheKeyCharacterInfo(characterId int32) string {
	return newCacheKey(
		pfx_nocorp,
		"CharacterInfo",
		strconv.Itoa(int(characterId)),
	)
}

func CacheKeyCorporationInfo(corporationId int32) string {
	return newCacheKey(
		pfx_nocorp,
		"CorporationInfo",
		strconv.Itoa(int(corporationId)),
	)
}

// structure info
func CacheKeyStructureInfo(structureId int64) string {
	return newCacheKey(
		pfx_corp,
		"StructureInfo",
		strconv.Itoa(int(structureId)),
	)
}

// // bucket

// web data
var CacheKeyWebBuybackSystemTypeMapsBuilder = newCacheKey(
	pfx_corp,
	b.OBJNAME_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
)
var CacheKeyWebBuybackBundleKeys = newCacheKey(
	CacheKeyWebBuybackSystemTypeMapsBuilder,
	"BundleKeys",
)
var CacheKeyWebShopLocationTypeMapsBuilder = newCacheKey(
	pfx_corp,
	b.OBJNAME_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
)
var CacheKeyWebShopBundleKeys = newCacheKey(
	CacheKeyWebShopLocationTypeMapsBuilder,
	"BundleKeys",
)
var CacheKeyWebBuybackSystems = newCacheKey(
	pfx_corp,
	b.OBJNAME_WEB_BUYBACK_SYSTEMS,
)
var CacheKeyWebShopLocations = newCacheKey(
	pfx_corp,
	b.OBJNAME_WEB_SHOP_LOCATIONS,
)
var CacheKeyWebMarkets = newCacheKey(
	pfx_corp,
	b.OBJNAME_WEB_MARKETS,
)

// build data
var CacheKeyBuildConstData = newCacheKey(
	pfx_corp,
	b.OBJNAME_CONSTANTS_DATA,
)

// auth hash set
func CacheKeyAuthHashSet(domain string) string {
	return newCacheKey(
		pfx_corp,
		"AuthHashSet",
		domain,
	)
}

// // RemoteDB

// user data
func CacheKeyNSUserData(characterId int32) string {
	return newCacheKey(
		pfx_corp,
		"UserData",
		strconv.Itoa(int(characterId)),
	)
}
func CacheKeyUserBuybackAppraisalCodes(characterId int32) string {
	cacheKeyUserData := CacheKeyNSUserData(characterId)
	return newCacheKey(
		cacheKeyUserData,
		"BuybackAppraisalCodes",
	)
}
func CacheKeyUserShopAppraisalCodes(characterId int32) string {
	cacheKeyUserData := CacheKeyNSUserData(characterId)
	return newCacheKey(
		cacheKeyUserData,
		"ShopAppraisalCodes",
	)
}
func CacheKeyUserCancelledPurchase(characterId int32) string {
	cacheKeyUserData := CacheKeyNSUserData(characterId)
	return newCacheKey(
		cacheKeyUserData,
		"CancelledPurchase",
	)
}
func CacheKeyUserMadePurchase(characterId int32) string {
	cacheKeyUserData := CacheKeyNSUserData(characterId)
	return newCacheKey(
		cacheKeyUserData,
		"MadePurchase",
	)
}

// appraisal
func CacheKeyAppraisal(appraisalCode string) string {
	return newCacheKey(
		pfx_corp,
		"Appraisal",
		appraisalCode,
	)
}

// purchase queue
var CacheKeyRawPurchaseQueue = newCacheKey(
	pfx_corp,
	"RawPurchaseQueue",
)

// // Composition
var CacheKeyPurchaseQueue = newCacheKey(
	pfx_corp,
	"PurchaseQueue",
)

func CacheKeyLocationPurchaseQueue(
	locationId int64,
) string {
	return newCacheKey(
		CacheKeyPurchaseQueue,
		strconv.Itoa(int(locationId)),
	)
}

var CacheKeyNSRawShopAssets = newCacheKey(
	pfx_corp,
	"RawShopAssets",
)

func CacheKeyRawShopAssets(locationId int64) string {
	return newCacheKey(
		CacheKeyNSRawShopAssets,
		strconv.Itoa(int(locationId)),
	)
}

var CacheKeyNSUnreservedShopAssets = newCacheKey(
	pfx_corp,
	"UnreservedShopAssets",
)

func CacheKeyUnreservedShopAssets(locationId int64) string {
	return newCacheKey(
		CacheKeyNSUnreservedShopAssets,
		strconv.Itoa(int(locationId)),
	)
}

var CacheKeyContracts = newCacheKey(
	pfx_corp,
	"Contracts",
)

// contractIds are unique
func CacheKeyContractItems(contractId int32) string {
	return newCacheKey(
		pfx_nocorp,
		"ContractItems",
		strconv.Itoa(int(contractId)),
	)
}

func CacheKeyNSRegionMarketOrders(
	regionId int32,
	typeId int32,
	isBuy bool,
) string {
	return newCacheKey(
		pfx_nocorp,
		"RegionMarketOrders",
		strconv.Itoa(int(regionId)),
		strconv.Itoa(int(typeId)),
		isBuyStr(isBuy),
	)
}

func CacheKeyRegionMarketOrders(
	nsCacheKey string,
	locationId int64,
) string {
	return newCacheKey(
		nsCacheKey,
		strconv.Itoa(int(locationId)),
	)
}

func CacheKeyNSStructureMarketOrders(structureId int64) string {
	return newCacheKey(
		pfx_nocorp,
		"StructureMarketOrders",
		strconv.Itoa(int(structureId)),
	)
}

func CacheKeyStructureMarketOrders(
	nsCacheKey string,
	typeId int32,
	isBuy bool,
) string {
	return newCacheKey(
		nsCacheKey,
		strconv.Itoa(int(typeId)),
		isBuyStr(isBuy),
	)
}

func CacheKeyTokenCharacter(
	app uint8,
	refreshToken string,
) string {
	return newCacheKey(
		pfx_corp,
		"TokenCharacter",
		strconv.Itoa(int(app)),
		refreshToken,
	)
}
