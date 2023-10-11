package cachekeys

import (
	"fmt"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WRITE_BUYBACK_APPRAISAL_NULL_ANTI_CACHE_KEY string = "WBA_NULL"

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

	IS_BUY_TRUE  string = "b"
	IS_BUY_FALSE string = "s"
)

func WriteBuybackAppraisalAntiCacheKey(characterId *int32) string {
	if characterId == nil {
		return WRITE_BUYBACK_APPRAISAL_NULL_ANTI_CACHE_KEY
	} else {
		return ReadUserDataCacheKey(*characterId)
	}
}

func AuthHashSetReaderCacheKey(authDomain string) string {
	return fmt.Sprintf("%s-%s", AUTH_HASH_SET_READER_PREFIX, authDomain)
}

func BucketReaderCacheKey(objName string) string {
	return objName
}

func ContractItemsCacheKey(contractId int32) string {
	return fmt.Sprintf(
		"%s-%d",
		CONTRACT_ITEMS_PREFIX,
		contractId,
	)
}

func ContractsCacheKey() string {
	return CONTRACTS_CACHE_KEY
}

func JWKSCacheKey() string {
	return JWKS_CACHE_KEY
}

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

func StructureInfoCacheKey(structureId int64) string {
	return fmt.Sprintf(
		"%s-%d",
		STRUCTURE_INFO_PREFIX,
		structureId,
	)
}

func AllShopAssetsCacheKey() string {
	return ALL_SHOP_ASSETS_CACHE_KEY
}

func UnreservedShopAssetsCacheKey() string {
	return UNRESERVED_SHOP_ASSETS_CACHE_KEY
}

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

func ReadAppraisalCacheKey(appraisalCode string) string {
	return appraisalCode
	// return fmt.Sprintf("appraisal-%s", appraisalCode)
}

func ReadUserDataCacheKey(characterId int32) string {
	return fmt.Sprintf(
		"%s-%d",
		READ_USER_DATA_PREFIX,
		characterId,
	)
}

func ReadShopQueueCacheKey() string {
	return SHOP_QUEUE_READ_CACHE_KEY
}

func WebBuybackSystemTypeMapsBuilderReaderCacheKey() string {
	return b.OBJNAME_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER
}

func WebBuybackBundleKeysCacheKey() string {
	return fmt.Sprintf(
		"%s-%s",
		WEB_BUYBACK_BUNDLE_KEYS_PREFIX,
		b.OBJNAME_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
	)
}

func WebShopBundleKeysCacheKey() string {
	return fmt.Sprintf(
		"%s-%s",
		WEB_SHOP_BUNDLE_KEYS_PREFIX,
		b.OBJNAME_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
	)
}

func WebShopLocationTypeMapsBuilderReaderCacheKey() string {
	return b.OBJNAME_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER
}

func WebBuybackSystemsReaderCacheKey() string {
	return b.OBJNAME_WEB_BUYBACK_SYSTEMS
}

func WebShopLocationsReaderCacheKey() string {
	return b.OBJNAME_WEB_SHOP_LOCATIONS
}

func WebMarketsReaderCacheKey() string {
	return b.OBJNAME_WEB_MARKETS
}

// // util
func isBuyStr(isBuy bool) string {
	if isBuy {
		return IS_BUY_TRUE
	} else {
		return IS_BUY_FALSE
	}
}
