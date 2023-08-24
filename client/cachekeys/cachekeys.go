package cachekeys

import (
	"fmt"

	cfg "github.com/WiggidyW/eve-trading-co-go/client/configure"
)

const (
	CONTRACTS_CACHE_KEY              = "contracts"
	JWKS_CACHE_KEY                   = "jwks"
	ALL_SHOP_ASSETS_CACHE_KEY        = "allassets"
	UNRESERVED_SHOP_ASSETS_CACHE_KEY = "unresassets"
	SHOP_QUEUE_READ_CACHE_KEY        = "shopqueue"

	REGION_MARKET_PREFIX                = "regionmarket"
	CHARACTER_INFO_PREFIX               = "charinfo"
	STRUCTURE_INFO_PREFIX               = "structureinfo"
	STRUCTURE_MARKET_PREFIX             = "structuremarket"
	READ_USER_DATA_PREFIX               = "chardata"
	RATE_LIMITING_CONTRACT_ITEMS_PREFIX = "contractitems"

	IS_BUY_TRUE  = "b"
	IS_BUY_FALSE = "s"
)

func BucketReaderCacheKey(objName string) string {
	return objName
}

func RateLimitingContractItemsCacheKey(contractId int32) string {
	return fmt.Sprintf(
		"%s-%d",
		RATE_LIMITING_CONTRACT_ITEMS_PREFIX,
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

func ShopQueueReadCacheKey() string {
	return SHOP_QUEUE_READ_CACHE_KEY
}

func GetBuybackSystemTypeMapsBuilderCacheKey() string {
	return cfg.B_TYPE_MAPS_BUILDER_DOMAIN_KEY
}

func GetShopLocationTypeMapsBuilderCacheKey() string {
	return cfg.S_TYPE_MAPS_BUILDER_DOMAIN_KEY
}

func GetBuybackSystemsCacheKey() string {
	return cfg.BUYBACK_SYSTEMS_DOMAIN_KEY
}

func GetShopLocationsCacheKey() string {
	return cfg.SHOP_LOCATIONS_DOMAIN_KEY
}

func GetMarketsCacheKey() string {
	return cfg.MARKETS_DOMAIN_KEY
}

// // util
func isBuyStr(isBuy bool) string {
	if isBuy {
		return IS_BUY_TRUE
	} else {
		return IS_BUY_FALSE
	}
}
