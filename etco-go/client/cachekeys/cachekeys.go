package cachekeys

import (
	"fmt"

	b "github.com/WiggidyW/etco-go-bucket"
	builder "github.com/WiggidyW/etco-go/buildconstants"
)

const (
	NULL_ANTI_CACHE_KEY string = "_ACNULL"

	CONTRACTS_CACHE_KEY              string = "contracts"
	JWKS_CACHE_KEY                   string = "jwks"
	ALL_SHOP_ASSETS_CACHE_KEY        string = "allassets"
	UNRESERVED_SHOP_ASSETS_CACHE_KEY string = "unresassets"
	SHOP_QUEUE_READ_CACHE_KEY        string = "shopqueue"

	CHARACTER_INFO_PREFIX          string = "charinfo"
	CORPORATION_INFO_PREFIX        string = "corpinfo"
	ALLIANCE_INFO_PREFIX           string = "alncinfo"
	REGION_MARKET_PREFIX           string = "regionmarket"
	STRUCTURE_INFO_PREFIX          string = "structureinfo"
	STRUCTURE_MARKET_PREFIX        string = "structuremarket"
	READ_USER_DATA_PREFIX          string = "chardata"
	CONTRACT_ITEMS_PREFIX          string = "contractitems"
	WEB_BUYBACK_BUNDLE_KEYS_PREFIX string = "bundlekeys"
	WEB_SHOP_BUNDLE_KEYS_PREFIX    string = "bundlekeys"
	AUTH_HASH_SET_READER_PREFIX    string = "authhashset"
	BUCKET_READER_PREFIX           string = "bucket"

	IS_BUY_TRUE  string = "b"
	IS_BUY_FALSE string = "s"
)

func WriteBuybackAppraisalAntiCacheKey(characterId *int32) string {
	if characterId == nil {
		return NULL_ANTI_CACHE_KEY
	} else {
		return ReadUserDataCacheKey(*characterId)
	}
}

// SHARED - JWKS is static
func JWKSCacheKey() string {
	return JWKS_CACHE_KEY
}

// SHARED - contractIds are unique
func ContractItemsCacheKey(contractId int32) string {
	return fmt.Sprintf(
		"%s-%d",
		CONTRACT_ITEMS_PREFIX,
		contractId,
	)
}

// SHARED - market results are the same for all corps
func RegionMarketCacheKey(
	regionId int32,
	typeId int32,
	isBuy bool,
) string {
	return fmt.Sprintf(
		"%s-%d-%d-%s",
		REGION_MARKET_PREFIX,
		regionId,
		typeId,
		isBuyStr(isBuy),
	)
}
func FilterRegionMarketCacheKey(
	regionId int32,
	typeId int32,
	isBuy bool,
	locationId int64,
) string {
	return fmt.Sprintf(
		"%s-%d",
		RegionMarketCacheKey(regionId, typeId, isBuy),
		locationId,
	)
}
func StructureMarketCacheKey(structureId int64) string {
	return fmt.Sprintf(
		"%s-%d",
		STRUCTURE_MARKET_PREFIX,
		structureId,
	)
}
func FilterStructureMarketCacheKey(
	structureId int64,
	typeId int32,
	isBuy bool,
) string {
	return fmt.Sprintf(
		"%s-%d-%s",
		StructureMarketCacheKey(structureId),
		typeId,
		isBuyStr(isBuy),
	)
}

// SHARED - entity info results are the same for all corps
func CharacterInfoCacheKey(characterId int32) string {
	return fmt.Sprintf(
		"%s-%d",
		CHARACTER_INFO_PREFIX,
		characterId,
	)
}
func CorporationInfoCacheKey(corporationId int32) string {
	return fmt.Sprintf(
		"%s-%d",
		CORPORATION_INFO_PREFIX,
		corporationId,
	)
}
func AllianceInfoCacheKey(allianceId int32) string {
	return fmt.Sprintf(
		"%s-%d",
		ALLIANCE_INFO_PREFIX,
		allianceId,
	)
}

// NOT SHARED - each corp has their own contracts
func ContractsCacheKey() string {
	return fmt.Sprintf(
		"%d-%s",
		builder.CORPORATION_ID,
		CONTRACTS_CACHE_KEY,
	)
}

// NOT SHARED - structure info depends upon docking rights, which are corp-specific
func StructureInfoCacheKey(structureId int64) string {
	return fmt.Sprintf(
		"%d-%s-%d",
		builder.CORPORATION_ID,
		STRUCTURE_INFO_PREFIX,
		structureId,
	)
}

// NOT SHARED - each corp has their own assets
func AllShopAssetsCacheKey() string {
	return fmt.Sprintf(
		"%d-%s",
		builder.CORPORATION_ID,
		ALL_SHOP_ASSETS_CACHE_KEY,
	)
}
func UnreservedShopAssetsCacheKey() string {
	return fmt.Sprintf(
		"%d-%s",
		builder.CORPORATION_ID,
		UNRESERVED_SHOP_ASSETS_CACHE_KEY,
	)
}

// NOT SHARED - corp-specific +
// negligible Sprintf performance cost reduces negligible hash collision odds
func ReadAppraisalCacheKey(appraisalCode string) string {
	return fmt.Sprintf(
		"%d-%s",
		builder.CORPORATION_ID,
		appraisalCode,
	)
}

// NOT SHARED - userData is corp-specific
func ReadUserDataCacheKey(characterId int32) string {
	return fmt.Sprintf(
		"%d-%s-%d",
		builder.CORPORATION_ID,
		READ_USER_DATA_PREFIX,
		characterId,
	)
}

// NOT SHARED - each corp has their own shop queue
func ReadShopQueueCacheKey() string {
	return fmt.Sprintf(
		"%d-%s",
		builder.CORPORATION_ID,
		SHOP_QUEUE_READ_CACHE_KEY,
	)
}

// NOT SHARED - bucket data is corp-specific
func AuthHashSetReaderCacheKey(authDomain string) string {
	return fmt.Sprintf(
		"%d-%s-%s",
		builder.CORPORATION_ID,
		AUTH_HASH_SET_READER_PREFIX,
		authDomain,
	)
}
func WebBuybackSystemTypeMapsBuilderReaderCacheKey() string {
	return fmt.Sprintf(
		"%d-%s",
		builder.CORPORATION_ID,
		b.OBJNAME_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
	)
}
func WebBuybackBundleKeysCacheKey() string {
	return fmt.Sprintf(
		"%d-%s-%s",
		builder.CORPORATION_ID,
		WEB_BUYBACK_BUNDLE_KEYS_PREFIX,
		b.OBJNAME_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
	)
}
func WebShopBundleKeysCacheKey() string {
	return fmt.Sprintf(
		"%d-%s-%s",
		builder.CORPORATION_ID,
		WEB_SHOP_BUNDLE_KEYS_PREFIX,
		b.OBJNAME_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
	)
}
func WebShopLocationTypeMapsBuilderReaderCacheKey() string {
	return fmt.Sprintf(
		"%d-%s",
		builder.CORPORATION_ID,
		b.OBJNAME_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
	)
}
func WebBuybackSystemsReaderCacheKey() string {
	return fmt.Sprintf(
		"%d-%s",
		builder.CORPORATION_ID,
		b.OBJNAME_WEB_BUYBACK_SYSTEMS,
	)
}
func WebShopLocationsReaderCacheKey() string {
	return fmt.Sprintf(
		"%d-%s",
		builder.CORPORATION_ID,
		b.OBJNAME_WEB_SHOP_LOCATIONS,
	)
}
func WebMarketsReaderCacheKey() string {
	return fmt.Sprintf(
		"%d-%s",
		builder.CORPORATION_ID,
		b.OBJNAME_WEB_MARKETS,
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
